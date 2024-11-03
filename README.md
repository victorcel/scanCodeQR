# Escáner de Códigos QR

Este proyecto es un servidor HTTP simple escrito en Go que escanea códigos QR de imágenes subidas y devuelve el contenido decodificado en formato JSON.

## Requisitos

- Go 1.23.2 o posterior
- Biblioteca `github.com/makiuchi-d/gozxing`
- Biblioteca `golang.org/x/image`

## Instalación

1. Clona el repositorio:
    ```sh
    git clone https://github.com/tuusuario/qrcode-scanner.git
    cd qrcode-scanner
    ```

2. Instala las dependencias:
    ```sh
    go mod tidy
    ```

## Uso

1. Compila el servidor para Linux con arquitectura amd64:
    ```sh
    GOOS=linux GOARCH=amd64 go build -o app
    ```

2. Ejecuta el servidor:
    ```sh
    ./main
    ```

3. Envía una solicitud POST a `http://localhost:8888/scan` con un cuerpo de datos de formulario que contenga un archivo de imagen con la clave `code`.

    Ejemplo usando `curl`:
    ```sh
    curl -X POST -F "code=@ruta/a/tu/imagen.png" http://localhost:8888/scan
    ```

## API

### POST /scan

#### Solicitud

- **Content-Type**: `multipart/form-data`
- **Datos del Formulario**:
  - `code`: El archivo de imagen que contiene el código QR.

#### Respuesta

- **Content-Type**: `application/json`
- **Cuerpo**:
  - En éxito: `{"result": "contenido decodificado del código QR"}`
  - En error: `{"error": "mensaje de error"}`

## Estructura del Proyecto

- `main.go`: El código principal del servidor.
- `go.mod`: El archivo del módulo Go.
- `.gitignore`: Archivo de git ignore.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT.