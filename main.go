package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type CameraServer struct {
	port    string
	running bool
}

type StreamInfo struct {
	Port       string    `json:"port"`
	Timestamp  time.Time `json:"timestamp"`
	Camera     string    `json:"camera"`
	Resolution string    `json:"resolution"`
}

func main() {
	server := &CameraServer{
		port: "8080",
	}

	http.HandleFunc("/", server.handleHome)
	http.HandleFunc("/stream", server.handleStream)
	http.HandleFunc("/api/status", server.handleStatus)
	http.HandleFunc("/api/start-camera", server.handleStartCamera)
	http.HandleFunc("/api/stop-camera", server.handleStopCamera)

	// Obtener IP local
	ip := getLocalIP()

	fmt.Printf("🎥 Alien Cam Server iniciado\n")
	fmt.Printf("📱 Acceso local: http://localhost:%s\n", server.port)
	fmt.Printf("💻 Acceso desde otros dispositivos: http://%s:%s\n", ip, server.port)
	fmt.Printf("🌐 Presiona Ctrl+C para detener\n\n")
	fmt.Printf("📋 Si la IP %s no funciona, intenta:\n", ip)
	fmt.Printf("   - Abrir Termux y ejecutar: ip route get 8.8.8.8\n")
	fmt.Printf("   - O revisar la configuración WiFi de tu celular\n\n")

	log.Fatal(http.ListenAndServe(":"+server.port, nil))
}

func (cs *CameraServer) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🎥 Alien Cam - Transmisión de Cámara</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            align-items: center;
            padding: 20px;
            color: white;
        }
        
        .container {
            width: 100%;
            max-width: 800px;
            background: rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(10px);
            border-radius: 20px;
            padding: 30px;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
        }
        
        h1 {
            text-align: center;
            margin-bottom: 30px;
            font-size: 2.5em;
            text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
        }
        
        .video-container {
            position: relative;
            width: 100%;
            background: #000;
            border-radius: 15px;
            overflow: hidden;
            margin-bottom: 20px;
            aspect-ratio: 16/9;
        }
        
        .video-placeholder {
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            background: linear-gradient(45deg, #1a1a2e, #16213e);
            color: #fff;
            font-size: 1.2em;
            text-align: center;
        }
        
        #videoStream {
            width: 100%;
            height: 100%;
            object-fit: cover;
            display: none;
        }
        
        .controls {
            display: flex;
            gap: 15px;
            justify-content: center;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }
        
        .btn {
            padding: 12px 24px;
            border: none;
            border-radius: 25px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        
        .btn-primary {
            background: linear-gradient(45deg, #00d4ff, #0099cc);
            color: white;
        }
        
        .btn-danger {
            background: linear-gradient(45deg, #ff416c, #ff4b2b);
            color: white;
        }
        
        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 20px rgba(0, 0, 0, 0.2);
        }
        
        .btn:disabled {
            opacity: 0.6;
            cursor: not-allowed;
            transform: none;
        }
        
        .status {
            text-align: center;
            padding: 15px;
            border-radius: 10px;
            margin-bottom: 20px;
            background: rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(5px);
        }
        
        .status-indicator {
            display: inline-block;
            width: 12px;
            height: 12px;
            border-radius: 50%;
            margin-right: 10px;
            background: #ff4444;
            animation: pulse 2s infinite;
        }
        
        .status-indicator.active {
            background: #44ff44;
        }
        
        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }
        
        .info {
            background: rgba(255, 255, 255, 0.05);
            padding: 20px;
            border-radius: 10px;
            margin-top: 20px;
        }
        
        .info h3 {
            margin-bottom: 10px;
            color: #00d4ff;
        }
        
        .info p {
            line-height: 1.6;
            margin-bottom: 10px;
        }
        
        .loading {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 3px solid rgba(255, 255, 255, 0.3);
            border-radius: 50%;
            border-top-color: white;
            animation: spin 1s ease-in-out infinite;
            margin-right: 10px;
        }
        
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
        
        @media (max-width: 600px) {
            .container {
                padding: 20px;
            }
            
            h1 {
                font-size: 2em;
            }
            
            .controls {
                flex-direction: column;
            }
            
            .btn {
                width: 100%;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎥 Alien Cam</h1>
        
        <div class="status">
            <span class="status-indicator" id="statusIndicator"></span>
            <span id="statusText">Cámara desactivada</span>
        </div>
        
        <div class="video-container">
            <div class="video-placeholder" id="placeholder">
                <div>
                    <p>📱 La cámara se mostrará aquí</p>
                    <p style="font-size: 0.8em; opacity: 0.7; margin-top: 10px;">Haz clic en "Iniciar Cámara" para comenzar</p>
                </div>
            </div>
            <img id="videoStream" alt="Transmisión de cámara">
        </div>
        
        <div class="controls">
            <button class="btn btn-primary" id="startBtn" onclick="startCamera()">
                🎥 Iniciar Cámara
            </button>
            <button class="btn btn-danger" id="stopBtn" onclick="stopCamera()" disabled>
                ⏹️ Detener Cámara
            </button>
        </div>
        
        <div class="info">
            <h3>📋 Información</h3>
            <p><strong>IP Local:</strong> <span id="localIP"></span></p>
            <p><strong>Puerto:</strong> 8080</p>
            <p><strong>Estado:</strong> <span id="detailedStatus">Esperando iniciar cámara</span></p>
            <p style="margin-top: 15px; font-size: 0.9em; opacity: 0.8;">
                💡 Para acceder desde otros dispositivos en la misma red, usa la IP local seguida del puerto 8080
            </p>
        </div>
    </div>

    <script>
        let isStreaming = false;
        
        // Obtener IP local
        fetch('/api/status')
            .then(response => response.json())
            .then(data => {
                document.getElementById('localIP').textContent = window.location.hostname;
            });
        
        function updateStatus(isActive, message) {
            const indicator = document.getElementById('statusIndicator');
            const statusText = document.getElementById('statusText');
            const detailedStatus = document.getElementById('detailedStatus');
            const startBtn = document.getElementById('startBtn');
            const stopBtn = document.getElementById('stopBtn');
            const placeholder = document.getElementById('placeholder');
            const videoStream = document.getElementById('videoStream');
            
            if (isActive) {
                indicator.classList.add('active');
                startBtn.disabled = true;
                stopBtn.disabled = false;
                placeholder.style.display = 'none';
                videoStream.style.display = 'block';
                videoStream.src = '/stream?' + new Date().getTime();
            } else {
                indicator.classList.remove('active');
                startBtn.disabled = false;
                stopBtn.disabled = true;
                placeholder.style.display = 'flex';
                videoStream.style.display = 'none';
                videoStream.src = '';
            }
            
            statusText.textContent = message;
            detailedStatus.textContent = message;
        }
        
        async function startCamera() {
            const startBtn = document.getElementById('startBtn');
            const originalText = startBtn.innerHTML;
            startBtn.innerHTML = '<span class="loading"></span>Iniciando...';
            startBtn.disabled = true;
            
            try {
                // Primero verificar si el endpoint de streaming está disponible
                const testStream = await fetch('/stream?' + new Date().getTime());
                
                if (testStream.ok) {
                    const response = await fetch('/api/start-camera', {
                        method: 'POST'
                    });
                    
                    if (response.ok) {
                        isStreaming = true;
                        updateStatus(true, 'Cámara activa y transmitiendo');
                        
                        // Cargar la primera imagen inmediatamente
                        const videoStream = document.getElementById('videoStream');
                        videoStream.src = '/stream?' + new Date().getTime();
                    } else {
                        throw new Error('Error al iniciar la cámara');
                    }
                } else {
                    // El streaming funciona pero con imagen de demostración
                    const response = await fetch('/api/start-camera', {
                        method: 'POST'
                    });
                    
                    if (response.ok) {
                        isStreaming = true;
                        updateStatus(true, 'Cámara demostración activa (instala Termux:API para acceso real)');
                        
                        const videoStream = document.getElementById('videoStream');
                        videoStream.src = '/stream?' + new Date().getTime();
                    } else {
                        throw new Error('Error al iniciar la cámara');
                    }
                }
            } catch (error) {
                console.error('Error:', error);
                updateStatus(false, 'Error al iniciar la cámara: ' + error.message);
                startBtn.innerHTML = originalText;
                startBtn.disabled = false;
            }
        }
        
        async function stopCamera() {
            const stopBtn = document.getElementById('stopBtn');
            const originalText = stopBtn.innerHTML;
            stopBtn.innerHTML = '<span class="loading"></span>Deteniendo...';
            stopBtn.disabled = true;
            
            try {
                const response = await fetch('/api/stop-camera', {
                    method: 'POST'
                });
                
                if (response.ok) {
                    isStreaming = false;
                    updateStatus(false, 'Cámara desactivada');
                } else {
                    throw new Error('Error al detener la cámara');
                }
            } catch (error) {
                console.error('Error:', error);
                stopBtn.disabled = false;
            }
            
            stopBtn.innerHTML = originalText;
        }
        
        // Actualizar stream periódicamente con manejo de errores
        setInterval(() => {
            if (isStreaming) {
                const videoStream = document.getElementById('videoStream');
                const newSrc = '/stream?' + new Date().getTime();
                
                // Verificar si la imagen anterior cargó correctamente
                videoStream.onerror = function() {
                    console.error('Error al cargar imagen del stream');
                    // Intentar recargar después de un pequeño delay
                    setTimeout(() => {
                        if (isStreaming) {
                            videoStream.src = '/stream?' + new Date().getTime();
                        }
                    }, 2000);
                };
                
                videoStream.src = newSrc;
            }
        }, 1000);
    </script>
</body>
</html>`

	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, nil)
}

func (cs *CameraServer) handleStream(w http.ResponseWriter, r *http.Request) {
	// Simular streaming de imagen desde cámara
	// En un entorno real, aquí se obtendría la imagen del dispositivo

	// Intentar capturar imagen usando Termux API
	imgData, err := cs.captureImage()
	if err != nil {
		// Si no hay Termux, servir una imagen de demostración
		w.Header().Set("Content-Type", "image/svg+xml")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `<svg width="640" height="480" xmlns="http://www.w3.org/2000/svg">
			<rect width="640" height="480" fill="%231a1a2e"/>
			<text x="320" y="240" font-family="Arial" font-size="24" fill="white" text-anchor="middle">
				📱 Cámara no disponible
			</text>
			<text x="320" y="270" font-family="Arial" font-size="16" fill="%23ccc" text-anchor="middle">
				Instala Termux:API para acceso real a la cámara
			</text>
		</svg>`)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusOK)
	w.Write(imgData)
}

func (cs *CameraServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	info := StreamInfo{
		Port:       cs.port,
		Timestamp:  time.Now(),
		Camera:     "Termux Camera",
		Resolution: "640x480",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func (cs *CameraServer) handleStartCamera(w http.ResponseWriter, r *http.Request) {
	cs.running = true
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "started",
		"message": "Cámara iniciada correctamente",
	})
}

func (cs *CameraServer) handleStopCamera(w http.ResponseWriter, r *http.Request) {
	cs.running = false
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "stopped",
		"message": "Cámara detendida correctamente",
	})
}

func (cs *CameraServer) captureImage() ([]byte, error) {
	// Siempre intentar con Termux API primero si estamos en Android
	if os.Getenv("TERMUX") != "" {
		// Crear directorio temporal si no existe
		tmpDir := "/data/data/com.termux/files/home"
		tmpFile := tmpDir + "/tmp_cam.jpg"

		// Verificar si Termux:API está disponible
		cmd := exec.Command("sh", "-c", "command -v termux-camera-info")
		if err := cmd.Run(); err == nil {
			// Capturar imagen con Termux:API
			cmd = exec.Command("termux-camera-photo", "-o", tmpFile)
			if err := cmd.Run(); err == nil {
				// Leer la imagen capturada
				imgData, err := os.ReadFile(tmpFile)
				if err == nil {
					return imgData, nil
				}
			}
		}
	}

	// Si no hay Termux o falla, retornar error para mostrar imagen demostración
	return nil, fmt.Errorf("camera not available")
}

func getLocalIP() string {
	// Método mejorado para obtener IP local en Termux
	if os.Getenv("TERMUX") != "" {
		// Método 1: ip route get
		cmd := exec.Command("sh", "-c", "ip route get 8.8.8.8 | awk '{print $7}' | head -1")
		if output, err := cmd.Output(); err == nil {
			ip := strings.TrimSpace(string(output))
			if ip != "" && ip != "0.0.0.0" {
				return ip
			}
		}

		// Método 2: hostname -I
		cmd = exec.Command("hostname", "-I")
		if output, err := cmd.Output(); err == nil {
			ips := strings.Fields(string(output))
			for _, ip := range ips {
				if strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") ||
					strings.HasPrefix(ip, "172.") || ip == "localhost" {
					if ip != "localhost" {
						return ip
					}
				}
			}
		}

		// Método 3: ifconfig
		cmd = exec.Command("sh", "-c", "ifconfig wlan0 2>/dev/null | grep 'inet ' | awk '{print $2}'")
		if output, err := cmd.Output(); err == nil {
			ip := strings.TrimSpace(string(output))
			if ip != "" && ip != "0.0.0.0" {
				return ip
			}
		}

		// Método 4: ip addr
		cmd = exec.Command("sh", "-c", "ip addr show wlan0 2>/dev/null | grep 'inet ' | awk '{print $2}' | cut -d'/' -f1")
		if output, err := cmd.Output(); err == nil {
			ip := strings.TrimSpace(string(output))
			if ip != "" && ip != "0.0.0.0" {
				return ip
			}
		}
	}

	// Fallback genérico
	return "localhost"
}
