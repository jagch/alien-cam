#!/bin/bash

# Script de compilación para Android/Termux
# Alien Cam - Compilación optimizada para Android

echo "🎥 Compilando Alien Cam para Android..."

# Verificar si estamos en Termux
if [ "$TERMUX" != "" ]; then
    echo "✅ Entorno Termux detectado"
    
    # Verificar dependencias
    echo "🔍 Verificando dependencias..."
    
    # Verificar Go
    if ! command -v go &> /dev/null; then
        echo "❌ Go no está instalado. Ejecuta: pkg install golang"
        exit 1
    fi
    
    # Verificar Termux:API
    if ! command -v termux-camera-info &> /dev/null; then
        echo "⚠️  Termux:API no detectado. La cámara funcionará en modo demostración"
        echo "   Para instalar: pkg install termux-api"
        echo "   Luego instala Termux:API desde Google Play y concede permisos"
    else
        echo "✅ Termux:API detectado"
    fi
    
    # Compilar para la arquitectura actual
    echo "🔨 Compilando para $(go env GOARCH)..."
    go build -o alien-cam main.go
    
    if [ $? -eq 0 ]; then
        echo "✅ Compilación exitosa"
        echo "📱 Ejecuta: ./alien-cam"
        echo "🌐 Accede a: http://localhost:8080"
    else
        echo "❌ Error en la compilación"
        exit 1
    fi
else
    echo "⚠️  Este script está diseñado para Termux/Android"
    echo "💻 Para compilar en otros sistemas, usa: go build -o alien-cam main.go"
fi