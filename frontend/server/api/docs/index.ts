import { readFileSync } from 'fs'
import { resolve } from 'path'
import { eventHandler, send } from 'h3'
import { fileURLToPath } from 'url'
import { join, dirname } from 'path'

// Get swagger spec
const swaggerJson = readFileSync(resolve('server/swagger/swagger.json'), 'utf8')

// Serve the Swagger HTML manually
export default eventHandler(async (event) => {
    const swaggerHtmlTemplate = `
    <!DOCTYPE html>
    <html>
    <head>
      <title>Swagger UI</title>
      <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui.css" />
    </head>
    <body>
      <div id="swagger-ui"></div>
      <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui-bundle.js"></script>
      <script>
        window.onload = () => {
          SwaggerUIBundle({
            spec: ${swaggerJson},
            dom_id: '#swagger-ui'
          })
        }
      </script>
    </body>
    </html>
  `

    return send(event, swaggerHtmlTemplate, 'text/html')
})