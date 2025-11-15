import json
from pathlib import Path
import pandas as pd
import joblib
from sentence_transformers import SentenceTransformer
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LogisticRegression
from sklearn.preprocessing import LabelEncoder
from sklearn.metrics import classification_report

# --- Configuration ---
DATA_PATH = "data/dataset.jsonl"
MODEL_DIR = Path("models")
EMBEDDING_MODEL_NAME = "all-MiniLM-L6-v2"
TEST_SIZE = 0.2
RANDOM_SEED = 42
RARE_CLASS_THRESHOLD = 2
MAX_ITERATIONS = 2000
# ---------------------

MODEL_DIR.mkdir(exist_ok=True)


def load_dataset() -> pd.DataFrame:
    """Loads the dataset from the JSONL file into a Pandas DataFrame."""
    rows = []
    with open(DATA_PATH, "r", encoding="utf-8") as f:
        for line in f:
            obj = json.loads(line)
            rows.append({"description": obj["description"], "label": obj["label"]})

    df = pd.DataFrame(rows)
    df["description"] = df["description"].astype(str).str.strip()
    return df


def preprocess_data(df: pd.DataFrame):
    """Splits, normalizes rare classes, and encodes labels."""

    # 1. Split Data First (Prevents Data Leakage)
    df_train, df_test = train_test_split(
        df, test_size=TEST_SIZE, random_state=RANDOM_SEED, stratify=df["label"]
    )
    df_train = df_train.reset_index(drop=True)
    df_test = df_test.reset_index(drop=True)

    # 2. Class Normalization (Fit only on Training)
    # Classes with < RARE_CLASS_THRESHOLD examples in the TRAINING set are converted to "other"
    train_counts = df_train["label"].value_counts()
    rare_classes = train_counts[train_counts < RARE_CLASS_THRESHOLD].index.tolist()

    # Apply normalization to training and test sets
    df_train.loc[df_train["label"].isin(rare_classes), "label"] = "other"
    df_test.loc[df_test["label"].isin(rare_classes), "label"] = "other"

    # 3. Label Encoding (Fit only on Training)
    le = LabelEncoder()
    y_train = le.fit_transform(df_train["label"])
    y_test = le.transform(df_test["label"])

    return df_train["description"], y_train, df_test["description"], y_test, le


def train_model(X_train, y_train, X_test, y_test, le: LabelEncoder):
    """Trains the Logistic Regression classifier and evaluates it."""

    # 5. Train and Evaluate Classifier
    clf = LogisticRegression(max_iter=MAX_ITERATIONS)
    clf.fit(X_train, y_train)

    preds = clf.predict(X_test)

    print("\n--- Classification Report ---")
    print(
        classification_report(
            y_test,
            preds,
            target_names=le.classes_,
            labels=list(range(len(le.classes_))),
        )
    )

    return clf


def save_models(clf, le, embedder: SentenceTransformer):
    """Saves the trained models (Classifier, Label Encoder, and Embedder)."""

    # a. Save the Classifier and LabelEncoder
    joblib.dump(clf, MODEL_DIR / "clf.joblib")
    joblib.dump(le, MODEL_DIR / "label_encoder.joblib")

    # b. Save the SentenceTransformer model
    embedder.save(str(MODEL_DIR / "embedder"))

    print(f"\n✅ Training complete. Models saved in the /{MODEL_DIR.name} folder.")


def main():
    """Main training pipeline."""
    df = load_dataset()
    print("Dataset rows:", len(df))

    # Preprocessing
    X_train_desc, y_train, X_test_desc, y_test, le = preprocess_data(df)

    # 4. Embeddings
    embedder = SentenceTransformer(EMBEDDING_MODEL_NAME)

    print("Generating Training embeddings...")
    X_train_embed = embedder.encode(
        X_train_desc.tolist(),
        show_progress_bar=True,
        convert_to_numpy=True,
    )

    print("Generating Test embeddings...")
    X_test_embed = embedder.encode(
        X_test_desc.tolist(),
        show_progress_bar=True,
        convert_to_numpy=True,
    )

    # Train and Save
    clf = train_model(X_train_embed, y_train, X_test_embed, y_test, le)
    save_models(clf, le, embedder)


if __name__ == "__main__":
    main()
