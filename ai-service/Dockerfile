FROM pytorch/pytorch:2.1.0-cuda11.8-cudnn8-runtime

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY main.py .

RUN python -c "from transformers import BlipProcessor, BlipForConditionalGeneration, CLIPProcessor, CLIPModel; \
    BlipProcessor.from_pretrained('Salesforce/blip-image-captioning-base'); \
    BlipForConditionalGeneration.from_pretrained('Salesforce/blip-image-captioning-base'); \
    CLIPProcessor.from_pretrained('openai/clip-vit-base-patch32'); \
    CLIPModel.from_pretrained('openai/clip-vit-base-patch32')"

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD python -c "import requests; requests.get('http://localhost:8000/health')" || exit 1

CMD ["python", "main.py"]
