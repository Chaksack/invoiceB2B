FROM n8nio/n8n

USER root
# For Alpine-based n8n image (default)
RUN apk add --no-cache poppler-utils tesseract-ocr tesseract-ocr-data-eng imagemagick
# poppler-utils provides pdftotext, pdftoppm, etc.
# tesseract-ocr for OCR
# imagemagick can be useful for image preprocessing before OCR

# If you were using a Debian-based n8n image:
# RUN apt-get update && \
#     apt-get install -y poppler-utils tesseract-ocr-eng imagemagick && \
#     rm -rf /var/lib/apt/lists/*
USER node