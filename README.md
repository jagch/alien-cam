# 🎥 Alien Cam - Transmisión de Cámara Android

Aplicación Go que transforma tu teléfono Android en una cámara IP accesible desde cualquier dispositivo en la misma red LAN.

## 📋 Requisitos

### En Android (Termux):
1. **Termux** - Emulador de terminal para Android
2. **Go** - Lenguaje de programación
3. **Termux:API** - Para acceso al hardware del dispositivo (opcional)

## 🚀 Instalación en Termux

### 1. Instalar Termux
Descarga Termux desde F-Droid: https://f-droid.org/packages/com.termux/

### 2. Actualizar paquetes
```bash
pkg update && pkg upgrade
```

### 3. Instalar Git y Go
```bash
pkg install git golang
```

### 4. Clonar y compilar Alien Cam
```bash
git clone <URL-del-repositorio>
cd alien-cam
go build -o alien-cam main.go
```

### 5. Ejecutar la aplicación
```bash
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
├── main.go          # Código principal del servidor
├── go.mod          # Módulo Go
├── README.md       # Este archivo
└── alien-cam       # Ejecutable compilado
```

## 🔒 Seguridad

- La aplicación solo escucha en la red local
- No almacena ni transmite datos externamente
- El streaming está limitado a la conexión actual

## 🤝 Contribuciones

¡Pull requests son bienvenidos!

## 📄 Licencia

MIT License