from fastapi import FastAPI
from pydantic import BaseModel
import joblib
from sentence_transformers import SentenceTransformer
import numpy as np
import uvicorn

# --- Configuration ---
MODEL_PATH_PREFIX = "models"
CONFIDENCE_THRESHOLD = 0.60
# ---------------------


# --- Pydantic Schema ---
class TextIn(BaseModel):
    """Input schema for text classification."""

    description: str


class ClassificationOut(BaseModel):
    """Output schema for classification result."""

    category: str
    confidence: float


# --- FastAPI App Initialization ---
app = FastAPI(
    title="Text Classification API",
    description="API for classifying text descriptions using Sentence Transformers and Logistic Regression.",
    version="1.0.0",
)

# --- Load Models ---
try:
    # Load scikit-learn models (Classifier and Label Encoder)
    clf = joblib.load(f"{MODEL_PATH_PREFIX}/clf.joblib")
    le = joblib.load(f"{MODEL_PATH_PREFIX}/label_encoder.joblib")

    # Load SentenceTransformer model
    embedder = SentenceTransformer(f"{MODEL_PATH_PREFIX}/embedder")
except FileNotFoundError as e:
    print(
        f"ERROR: Could not find model files. Please ensure the 'models' directory is present and contains all necessary files. Error: {e}"
    )
    # Optional: Exit or raise exception if models are critical
except Exception as e:
    print(f"An unexpected error occurred during model loading: {e}")
    # Optional: Exit or raise exception


@app.post(
    "/classify",
    response_model=ClassificationOut,
    summary="Classify a text description",
    description="Predicts the category of a given text description. Returns 'unknown' if confidence is below the set threshold.",
)
def classify_text(input_data: TextIn):
    """
    Classifies the input text description.
    """
    text = input_data.description

    # 1. Generate embedding
    emb = embedder.encode([text], convert_to_numpy=True)

    # 2. Predict probabilities
    # Predict_proba returns [[prob_0, prob_1, ..., prob_n]], so we take the first element [0]
    probs = clf.predict_proba(emb)[0]

    # 3. Find the best prediction
    idx = int(np.argmax(probs))
    confidence = float(probs[idx])

    # 4. Convert numerical index back to category name
    label = le.inverse_transform([idx])[0]

    # 5. Apply confidence threshold
    if confidence < CONFIDENCE_THRESHOLD:
        # Return 'unknown' if confidence is too low
        return ClassificationOut(category="unknown", confidence=confidence)

    # 6. Return the predicted category
    return ClassificationOut(category=label, confidence=confidence)


if __name__ == "__main__":
    # Note: uvicorn is imported at the top now for cleaner structure
    # The host is usually '127.0.0.1' or 'localhost' for local development,
    # but '0.0.0.0' is kept for maximum compatibility (e.g., Docker)
    uvicorn.run(app, host="0.0.0.0", port=8000)
