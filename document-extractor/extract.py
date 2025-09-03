# extractor.py
import pdfplumber
import sys


def extract_text(pdf_path, password):
    all_text = []

    with pdfplumber.open(pdf_path, password=password) as pdf:
        for page in pdf.pages:
            text = page.extract_text()
            if text:
                all_text.append(text)

    return "\n".join(all_text)


if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python extractor.py archivo.pdf password")
        sys.exit(1)

    pdf_file = sys.argv[1]
    password = sys.argv[2]

    try:
        text = extract_text(pdf_file, password)

        print(text)

    except Exception as e:
        print(e)
