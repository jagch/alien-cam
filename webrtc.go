package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type WebRTCManager struct {
	peerConnections map[string]*webrtc.PeerConnection
	mutex           sync.RWMutex
	upgrader        websocket.Upgrader
}

type SignalingMessage struct {
	Type    string      `json:"type"`
	PeerID  string      `json:"peerId"`
	Payload interface{} `json:"payload"`
}

func NewWebRTCManager() *WebRTCManager {
	return &WebRTCManager{
		peerConnections: make(map[string]*webrtc.PeerConnection),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Permitir todos los orígenes para desarrollo
			},
		},
	}
}

func (w *WebRTCManager) createPeerConnection(peerID string) (*webrtc.PeerConnection, error) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	// Configurar para recibir video
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("📥 Track recibido: %s", track.Codec().MimeType)

		// Aquí se procesaría el video de la cámara
		// Por ahora, solo logueamos que recibimos el track
		rtpBuf := make([]byte, 1500)
		for {
			_, _, readErr := track.Read(rtpBuf)
			if readErr != nil {
				log.Printf("❌ Error leyendo track: %v", readErr)
				return
			}
		}
	})

	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			log.Printf("✅ ICE gathering completado para peer %s", peerID)
			return
		}
		log.Printf("🧊 ICE candidato para peer %s: %s", peerID, candidate.String())
	})

	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Printf("🔄 Estado de conexión peer %s: %s", peerID, state.String())
		if state == webrtc.PeerConnectionStateFailed || state == webrtc.PeerConnectionStateClosed {
			w.removePeerConnection(peerID)
		}
	})

	w.mutex.Lock()
	w.peerConnections[peerID] = peerConnection
	w.mutex.Unlock()

	return peerConnection, nil
}

func (w *WebRTCManager) removePeerConnection(peerID string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if pc, exists := w.peerConnections[peerID]; exists {
		pc.Close()
		delete(w.peerConnections, peerID)
		log.Printf("🗑️  Peer connection %s eliminada", peerID)
	}
}

func (w *WebRTCManager) handleWebSocket(c *gin.Context) {
	conn, err := w.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ Error WebSocket upgrade: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("🔌 Cliente WebSocket conectado")

	for {
		var msg SignalingMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("❌ Error leyendo mensaje WebSocket: %v", err)
			break
		}

		log.Printf("📨 Mensaje recibido: %s de peer %s", msg.Type, msg.PeerID)

		switch msg.Type {
		case "offer":
			w.handleOffer(conn, msg)
		case "answer":
			w.handleAnswer(conn, msg)
		case "ice-candidate":
			w.handleICECandidate(conn, msg)
		}
	}
}

func (w *WebRTCManager) handleOffer(conn *websocket.Conn, msg SignalingMessage) {
	peerID := msg.PeerID

	pc, err := w.createPeerConnection(peerID)
	if err != nil {
		log.Printf("❌ Error creando peer connection: %v", err)
		return
	}

	offerData, _ := json.Marshal(msg.Payload)
	offer := webrtc.SessionDescription{}
	if err := json.Unmarshal(offerData, &offer); err != nil {
		log.Printf("❌ Error parseando offer: %v", err)
		return
	}

	if err := pc.SetRemoteDescription(offer); err != nil {
		log.Printf("❌ Error estableciendo remote description: %v", err)
		return
	}

	// Crear answer
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		log.Printf("❌ Error creando answer: %v", err)
		return
	}

	if err := pc.SetLocalDescription(answer); err != nil {
		log.Printf("❌ Error estableciendo local description: %v", err)
		return
	}

	// Enviar answer al cliente
	response := SignalingMessage{
		Type:    "answer",
		PeerID:  peerID,
		Payload: answer,
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("❌ Error enviando answer: %v", err)
	}
}

func (w *WebRTCManager) handleAnswer(conn *websocket.Conn, msg SignalingMessage) {
	// Implementar si el servidor inicia la conexión
	log.Printf("📋 Answer recibido para peer %s", msg.PeerID)
}

func (w *WebRTCManager) handleICECandidate(conn *websocket.Conn, msg SignalingMessage) {
	w.mutex.RLock()
	pc, exists := w.peerConnections[msg.PeerID]
	w.mutex.RUnlock()

	if !exists {
		log.Printf("❌ Peer connection %s no encontrada", msg.PeerID)
		return
	}

	candidateData, _ := json.Marshal(msg.Payload)
	candidate := webrtc.ICECandidateInit{}
	if err := json.Unmarshal(candidateData, &candidate); err != nil {
		log.Printf("❌ Error parseando ICE candidate: %v", err)
		return
	}

	if err := pc.AddICECandidate(candidate); err != nil {
		log.Printf("❌ Error añadiendo ICE candidate: %v", err)
	}
}

func (w *WebRTCManager) startVideoCapture(peerID string) {
	// Aquí implementaríamos la captura de video real
	// Por ahora, simulamos con un track de video
	log.Printf("🎥 Iniciando captura de video para peer %s", peerID)

	// En una implementación real:
	// 1. Usar FFmpeg para capturar video de la cámara Android
	// 2. Codificar a formato WebRTC (VP8/H264)
	// 3. Enviar paquetes RTP al peer connection
}
