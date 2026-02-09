# 🎥 Alien Cam - Transmisión de Cámara Android

Aplicación Go optimizada para Android que transforma tu teléfono en una cámara IP accesible desde cualquier dispositivo en la misma red LAN.

## 📋 Requisitos del Sistema

### Android:
- **Android 5.0+** (API 21)
- **Termux** - Emulador de terminal para Android
- **Go 1.21+** - Lenguaje de programación
- **Termux:API** - Para acceso real a la cámara (opcional)

### Permisos necesarios:
- Internet (para servidor web)
- Cámara (si se usa Termux:API)

## 🚀 Instalación en Android (Termux)

### 1. Instalar Termux
Descarga Termux desde **F-Droid** (recomendado): https://f-droid.org/packages/com.termux/

### 2. Actualizar paquetes
```bash
pkg update && pkg upgrade
```

### 3. Instalar dependencias
```bash
pkg install git golang
```

### 4. Instalar Termux:API (opcional, para cámara real)
```bash
# Instalar el paquete
pkg install termux-api

# Descargar Termux:API desde F-Droid o Google Play
# Conceder permisos de cámara cuando se solicite
```

### 5. Compilar y ejecutar
```bash
# Clonar o descargar el código
git clone <repository-url>
cd alien-cam

# Hacer ejecutable el script de compilación
chmod +x build-android.sh

# Compilar
./build-android.sh

# Ejecutar
./alien-cam
```

## 📱 Configuración de Cámara (Opcional pero Recomendado)

Para acceso real a la cámara del dispositivo:

1. Instalar **Termux:API** desde F-Droid
2. Conceder permisos de cámara a Termux
3. Ejecutar la aplicación con acceso a cámara

```bash
# Probar acceso a cámara
termux-camera-info
```

## 🌐 Acceso Web

Una vez iniciada la aplicación, verás algo como:
```
🎥 Alien Cam Server iniciado
📱 Acceso local: http://192.168.1.100:8080
💻 Acceso desde otros dispositivos: http://192.168.1.100:8080
```

### Desde el mismo dispositivo:
- Abre el navegador y visita `http://localhost:8080`

### Desde otros dispositivos en la misma red:
- Reemplaza con la IP que muestra la aplicación
- Ejemplo: `http://192.168.1.100:8080`

## ⚙️ Funcionalidades

- **Streaming en tiempo real** de la cámara del dispositivo
- **Interfaz web moderna** con controles intuitivos
- **Acceso multi-dispositivo** desde cualquier navegador
- **Indicadores de estado** en tiempo real
- **Diseño responsive** para móviles y escritorio

## 🔧 Uso

1. **Iniciar la aplicación**: Ejecuta `./alien-cam`
2. **Abrir navegador**: Ingresa la URL mostrada
3. **Iniciar cámara**: Haz clic en "Iniciar Cámara"
4. **Ver streaming**: La imagen aparecerá en la página web
5. **Compartir acceso**: Otros dispositivos pueden ver usando la misma IP

## 🛠️ Solución de Problemas

### Permiso denegado:
```bash
chmod +x alien-cam
```

### Puerto en uso:
El servidor usa el puerto 8080. Si está ocupado, cambia el puerto en el código.

### Cámara no disponible:
- Instala Termux:API
- Concede permisos de cámara en Android Settings
- Reinicia Termux

### Acceso desde otros dispositivos falla:
- Verifica que ambos dispositivos estén en la misma red WiFi
- Confirma que el firewall no bloquee el puerto 8080
- Usa la IP correcta que muestra la aplicación

## 📦 Estructura del Proyecto

```
alien-cam/
├── main.go              # Código principal del servidor
├── go.mod              # Módulo Go
├── build-android.sh    # Script de compilación para Android
├── README.md           # Este archivo
└── alien-cam           # Ejecutable compilado
```

## 🔧 Compilación Manual

Si el script automático no funciona:
```bash
# Verificar dependencias
go version

# Compilar manualmente
go build -o alien-cam main.go

# Ejecutar
./alien-cam
```

## 🔒 Seguridad

- La aplicación solo escucha en la red local
- No almacena ni transmite datos externamente
- El streaming está limitado a la conexión actual

## 🤝 Contribuciones

¡Pull requests son bienvenidos!

## 📄 Licencia

MIT License