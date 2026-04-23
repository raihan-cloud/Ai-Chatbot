package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	r := gin.Default()

	// Setup CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Ambil API Key sekali saat server start
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		panic("[ERROR] GEMINI_API_KEY tidak ditemukan di environment variable")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		panic(fmt.Sprintf("[ERROR] Gagal membuat client: %v", err))
	}
	defer client.Close()

	r.POST("/chat", func(c *gin.Context) {
		var input struct {
			Message string `json:"message"`
		}

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
			return
		}

		// Gunakan model flash-lite untuk kecepatan tinggi
		model := client.GenerativeModel("gemini-2.5-flash-lite")

		// --- BAGIAN KONFIGURASI CASE CLOUD & IOT ---
		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{
				genai.Text("Kamu adalah pakar Cloud Computing dan IoT. " +
					"Tugasmu adalah menjawab pertanyaan seputar DevOps, AWS, Azure, Google Cloud, ESP32, MQTT, dan sensor. " +
					"ATURAN KERAS: Jangan gunakan format Markdown. Jangan gunakan simbol **, ##, list -, atau tabel. " +
					"Berikan jawaban dalam bentuk teks polos (plain text). " +
					"Jika user bertanya di luar Cloud atau IoT, arahkan mereka kembali ke topik tersebut dengan sopan."),
			},
		}

		resp, err := model.GenerateContent(ctx, genai.Text(input.Message))
		if err != nil {
			fmt.Printf("[ERROR] Gemini API Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "AI gagal merespon: " + err.Error()})
			return
		}

		var reply string
		if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
			for _, part := range resp.Candidates[0].Content.Parts {
				reply += fmt.Sprintf("%v", part)
			}
		} else {
			reply = "Maaf, AI tidak memberikan jawaban."
		}

		c.JSON(http.StatusOK, gin.H{"reply": reply})
	})

	fmt.Println("Server backend berjalan di port :8080...")
	r.Run(":8080")
}