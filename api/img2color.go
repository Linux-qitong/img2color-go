package handler

import (
    "crypto/md5"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "image"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"

    "github.com/disintegration/imaging"
    "github.com/go-redis/redis/v8"
    "github.com/joho/godotenv"
    "github.com/lucasb-eyer/go-colorful"
    "github.com/nfnt/resize"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "golang.org/x/image/webp" // webp格式
    "golang.org/x/net/context"
    "github.com/gen2brain/avif" // 新增：avif格式
)

func extractMainColor(imgURL string) (string, error) {
    md5Hash := calculateMD5Hash([]byte(imgURL))

    if cacheEnabled && redisClient != nil {
        cachedColor, err := redisClient.Get(ctx, md5Hash).Result()
        if err == nil && cachedColor != "" {
            return cachedColor, nil
        }
    }

    req, err := http.NewRequest("GET", imgURL, nil)
    if err != nil {
        return "", err
    }

    req.Header.Set("User-Agent", "Mozilla/5.0")

    client := http.DefaultClient
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var img image.Image

    contentType := resp.Header.Get("Content-Type")
    switch contentType {
    case "image/webp":
        img, err = webp.Decode(resp.Body)
    case "image/avif": // 新增：avif格式处理
        img, err = avif.Decode(resp.Body)
    default:
        img, err = imaging.Decode(resp.Body)
    }

    if err != nil {
        return "", err
    }

    img = resize.Resize(50, 0, img, resize.Lanczos3)

    bounds := img.Bounds()
    var r, g, b uint32
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            c := img.At(x, y)
            r0, g0, b0, _ := c.RGBA()
            r += r0
            g += g0
            b += b0
        }
    }

    totalPixels := uint32(bounds.Dx() * bounds.Dy())
    averageR := r / totalPixels
    averageG := g / totalPixels
    averageB := b / totalPixels

    mainColor := colorful.Color{R: float64(averageR) / 0xFFFF, G: float64(averageG) / 0xFFFF, B: float64(averageB) / 0xFFFF}

    colorHex := mainColor.Hex()

    if cacheEnabled && redisClient != nil {
        _, err := redisClient.Set(ctx, md5Hash, colorHex, 0).Result()
        if err != nil {
            log.Printf("将结果存储在缓存中时出错：%v\n", err)
        }
    }

    if useMongoDB && colorsCollection != nil {
        _, err := colorsCollection.InsertOne(ctx, bson.M{
            "url":   imgURL,
            "color": colorHex,
        })
        if err != nil {
            log.Printf("将结果存储在MongoDB中时出错：%v\n", err)
        }
    }

    return colorHex, nil
}package handler

import (
    "crypto/md5"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "image"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"

    "github.com/disintegration/imaging"
    "github.com/go-redis/redis/v8"
    "github.com/joho/godotenv"
    "github.com/lucasb-eyer/go-colorful"
    "github.com/nfnt/resize"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "golang.org/x/image/webp" // webp格式
    "golang.org/x/net/context"
    "github.com/gen2brain/avif" // 新增：avif格式
)

func extractMainColor(imgURL string) (string, error) {
    md5Hash := calculateMD5Hash([]byte(imgURL))

    if cacheEnabled && redisClient != nil {
        cachedColor, err := redisClient.Get(ctx, md5Hash).Result()
        if err == nil && cachedColor != "" {
            return cachedColor, nil
        }
    }

    req, err := http.NewRequest("GET", imgURL, nil)
    if err != nil {
        return "", err
    }

    req.Header.Set("User-Agent", "Mozilla/5.0")

    client := http.DefaultClient
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var img image.Image

    contentType := resp.Header.Get("Content-Type")
    switch contentType {
    case "image/webp":
        img, err = webp.Decode(resp.Body)
    case "image/avif": // 新增：avif格式处理
        img, err = avif.Decode(resp.Body)
    default:
        img, err = imaging.Decode(resp.Body)
    }

    if err != nil {
        return "", err
    }

    img = resize.Resize(50, 0, img, resize.Lanczos3)

    bounds := img.Bounds()
    var r, g, b uint32
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            c := img.At(x, y)
            r0, g0, b0, _ := c.RGBA()
            r += r0
            g += g0
            b += b0
        }
    }

    totalPixels := uint32(bounds.Dx() * bounds.Dy())
    averageR := r / totalPixels
    averageG := g / totalPixels
    averageB := b / totalPixels

    mainColor := colorful.Color{R: float64(averageR) / 0xFFFF, G: float64(averageG) / 0xFFFF, B: float64(averageB) / 0xFFFF}

    colorHex := mainColor.Hex()

    if cacheEnabled && redisClient != nil {
        _, err := redisClient.Set(ctx, md5Hash, colorHex, 0).Result()
        if err != nil {
            log.Printf("将结果存储在缓存中时出错：%v\n", err)
        }
    }

    if useMongoDB && colorsCollection != nil {
        _, err := colorsCollection.InsertOne(ctx, bson.M{
            "url":   imgURL,
            "color": colorHex,
        })
        if err != nil {
            log.Printf("将结果存储在MongoDB中时出错：%v\n", err)
        }
    }

    return colorHex, nil
}
