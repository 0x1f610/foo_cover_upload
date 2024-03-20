package api

import (
	"bytes"
	"context"
	"crypto/md5"
	"embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/image/draw"

	"github.com/0x1f610/foo_cover_upload/utils"
)

var (
	ctx = context.Background()
)

func RedisClient(redisHost string, redisPassword string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("ERROR: Cannot connect to redis server.")
		fmt.Println(err)
		os.Exit(1)
	}

	return client
}

func Run(FS embed.FS, host string, address string, passHash string, redisHost string, redisPassword string) {

	gin.SetMode(gin.ReleaseMode)

	fmt.Println("Attempting to connect to redis server at " + redisHost)

	rdb := RedisClient(redisHost, redisPassword)

	fmt.Println("Connected to redis server at " + redisHost)

	r := gin.Default()
	templ := template.Must(template.New("").ParseFS(FS, "html/*.tmpl"))
	r.SetHTMLTemplate(templ)

	fe, _ := fs.Sub(FS, "html/static")
	r.StaticFS("/static", http.FS(fe))

	r.POST("/upload", func(c *gin.Context) {
		if passHash != "" {
			passAttempt := c.Request.Header.Get("Authorization")

			hasher := md5.New()
			hasher.Write([]byte(passAttempt))

			passAttemptHash := hasher.Sum(nil)

			if hex.EncodeToString(passAttemptHash[:]) != passHash {
				c.String(403, "Authentication required.")
				return
			}
		}

		data, err := io.ReadAll(c.Request.Body)

		if err != nil {
			c.String(400, "There is an error processing your request.\n\n"+err.Error())
			return
		}

		if len(data) > 10000000 {
			c.String(413, "Image larger than 10MB is not allowed.")
			return
		}

		dataContentType := http.DetectContentType(data)

		if dataContentType == "image/png" || dataContentType == "image/jpeg" {
			// Image is ok

			// Check if DB has room for storage
			// Hardcode limit is 10k, which takes about 550MB of RAM.
			// There really is no point going further, but I can make it configurable later.
			entries, err := rdb.DBSize(ctx).Result()
			if err != nil {
				c.String(500, "Database error.")
				return
			} else if entries >= 10000 {
				c.String(507, "Database full, try again later.")
				return
			}

			imageBytes := bytes.NewReader(data)

			imageData, _, err := image.Decode(imageBytes)

			if err != nil {
				c.String(400, "There is an error processing your request.\n\n"+err.Error())
				return
			}

			imageBounds := imageData.Bounds()
			imageWidth := imageBounds.Dx()
			imageHeight := imageBounds.Dy()

			// Reject non-square images
			if imageWidth != imageHeight {
				c.String(400, "Please upload only square images.")
				return
			}

			// Reject unreasonably big images
			if imageWidth > 5000 {
				c.String(413, "Your image is too big.")
				return
			}

			var imageToSave []byte

			// Only resize if image size is over 512
			if imageWidth > 512 {
				finalImage := image.NewRGBA(image.Rect(0, 0, 512, 512))

				draw.BiLinear.Scale(finalImage, finalImage.Bounds(), imageData, imageData.Bounds(), draw.Over, nil)

				buf := new(bytes.Buffer)
				err = png.Encode(buf, finalImage)
				if err != nil {
					panic(err)
				} else {
					imageToSave = buf.Bytes()
				}

			} else {
				imageToSave = data
			}

			// Generate a random string, and test if it exists
			// If not exists, write to redis with the key and data
			var imageKey string

			for {
				imageKey = utils.GenerateString(8)
				_, err := rdb.Get(ctx, imageKey).Result()
				if err == redis.Nil {
					break // Key doesn't exist, exit the loop
				} else if err != nil {
					// If it's other error, throw
					panic(err)
				}
			}

			// fmt.Println(len(data))
			// fmt.Println(len(imageToSave))

			// Write to redis - TTL 30 minutes
			// Half an hour should be enough for most of the albums out there
			err = rdb.Set(ctx, imageKey, imageToSave, 30*time.Minute).Err()
			if err != nil {
				panic(err)
			}
			c.String(200, host+"/image/"+imageKey)
			return
		} else {
			c.String(415, "Your image is not accepted.")
			return
		}
	})
	r.GET("/image/:key", func(c *gin.Context) {
		key := c.Param("key")
		if len(key) == 0 || len(key) > 8 {
			c.String(400, "Invalid key")
			return
		}
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				c.String(404, "Not found")
				return
			} else {
				c.String(500, "Internal server error")
				panic(err)
			}
		}
		c.Header("Content-Type", "image/png")
		c.Writer.WriteString(val)
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{
			"host": host,
		})
	})
	r.NoRoute(func(c *gin.Context) {
		c.String(404, "Not found")
	})

	fmt.Println("Image server started on " + address)
	fmt.Println()

	r.Run(address)
}
