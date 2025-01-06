# Rest API Go Events

## Descripción

Esta es una API RESTful construida con Go que permite gestionar eventos, usuarios y registros. La aplicación permite crear, leer, actualizar y eliminar eventos, así como gestionar registros de usuarios para eventos específicos.

## Características

- **Gestión de Eventos**: Crear, obtener, actualizar y eliminar eventos.
- **Registro de Usuarios**: Permite a los usuarios registrarse y gestionar su información.
- **Búsqueda de Eventos**: Buscar eventos por nombre, categoría, fecha y etiquetas.
- **Paginación**: Obtener resúmenes de eventos con paginación.
- **Integración con API de Clima**: Obtener información del clima para eventos próximos.
- **Envío de Emails**: Envía emails para recuperar y restablecer contraseñas.

## Tecnologías Utilizadas

- Go
- PostgreSQL
- Gin (framework web)
- Docker
- GitHub Actions (para CI/CD)

## Requisitos

- Go 1.16 o superior
- PostgreSQL
- Docker (opcional, para ejecutar en contenedores)

## Instalación

1. Clona el repositorio:

   ```bash
   git clone https://github.com/tu_usuario/restApi-Go.git
   cd restApi-Go
   ```

2. Configura las variables de entorno en un archivo `.env`:

   ```plaintext
   EMAIL_FROM= tu_email
   EMAIL_PASSWORD= tu_contraseña
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=tu_usuario
   DB_PASSWORD=tu_contraseña
   DB_NAME=tu_base_de_datos
   WEATHER_API_KEY=tu_api_key
   ```

3. Instala las dependencias:

   ```bash
   go mod tidy
   ```

4. Inicia la base de datos y crea las tablas necesarias:

   ```bash
   go run cmd/api/main.go
   ```

## Uso

### Endpoints

#### 🌍 Públicos

- **GET /events**: Obtener todos los eventos.
- **GET /events/:id**: Obtener un evento por ID.
- **GET /events/by-name**: Buscar eventos por nombre.
- **GET /events/by-tags**: Buscar eventos por etiquetas.
- **GET /events/by-category**: Buscar eventos por categoría.
- **GET /events/summaries**: Obtener resúmenes de eventos con paginación.
- **GET /tags**: Obtener todas las etiquetas.
- **GET /events/categories**: Obtener todas las categorías.

#### 🔒 Privados (requieren autenticación)

- **POST /events**: Crear un nuevo evento.
- **PUT /events/:id**: Actualizar un evento existente.
- **DELETE /events/:id**: Eliminar un evento.
- **POST /events/:id/register**: Registrar a un usuario en un evento.
- **DELETE /events/:id/register**: Cancelar la inscripción de un usuario en un evento.
- **GET /users/:id**: Obtener información de un usuario por ID.
- **PUT /users/:id**: Actualizar información de un usuario.
- **DELETE /users/:id**: Eliminar un usuario.

### Ejemplo de Solicitud

Para obtener todos los eventos:

```bash
curl -X GET http://localhost:8080/events
```

Para buscar eventos por nombre:

```bash
curl -X GET "http://localhost:8080/events/by-name?name=torneo"
```

## Contribuciones

Las contribuciones son bienvenidas. Si deseas contribuir, por favor sigue estos pasos:

1. Haz un fork del repositorio.
2. Crea una nueva rama (`git checkout -b feature/nueva-caracteristica`).
3. Realiza tus cambios y haz commit (`git commit -m 'Agrega nueva característica'`).
4. Haz push a la rama (`git push origin feature/nueva-caracteristica`).
5. Abre un Pull Request.

## Licencia

Este proyecto está bajo la Licencia MIT. Consulta el archivo [LICENSE](LICENSE) para más detalles.

## Contacto

Si tienes preguntas o sugerencias, no dudes en contactarme a través de [agustin.molina.dev@gmail.com](agustin.molina.dev@gmail.com).
