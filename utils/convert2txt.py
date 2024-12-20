import os
from PyPDF2 import PdfReader
from docx import Document


def convert_pdf_to_txt(pdf_path, txt_path):
    """Convert a PDF file to a TXT file."""
    print(f"Try to convert PDF to TXT: {pdf_path}")
    try:
        reader = PdfReader(pdf_path)
        with open(txt_path, "w", encoding="utf-8") as txt_file:
            for page in reader.pages:
                txt_file.write(page.extract_text() or "")
        print(f"Success.")
    except Exception as e:
        print(f"Failed, error: {e}")


def convert_docx_to_txt(docx_path, txt_path):
    """Convert a DOCX file to a TXT file."""
    print(f"Try to convert DOCX to TXT: {docx_path}")
    try:
        # Load the .docx file
        doc = Document(docx_path)
        # Extract all text content
        text = "\n".join([paragraph.text for paragraph in doc.paragraphs])
        # Write the text to the .txt file
        with open(txt_path, "w", encoding="utf-8") as txt_file:
            txt_file.write(text)
        print(f"Success.")
    except Exception as e:
        print(f"Failed, error: {e}")


def process_folder(input_folder, output_folder):
    """Convert all PDF and DOC/DOCX files in the folder to TXT."""
    if not os.path.exists(output_folder):
        os.makedirs(output_folder)

    for root, _, files in os.walk(input_folder):
        for file in files:
            file_path = os.path.join(root, file)
            file_name, file_ext = os.path.splitext(file)
            output_file = os.path.join(output_folder, f"{file_name}.txt")

            if file_ext.lower() == ".pdf":
                convert_pdf_to_txt(file_path, output_file)
            elif file_ext.lower() == ".docx":
                convert_docx_to_txt(file_path, output_file)
            elif file_ext.lower() == ".doc":
                print(f"Try to convert DOC to TXT: {file_path}")
                print(f"Unsupported file type: .doc")


if __name__ == "__main__":
    input_folder = input("Enter the path to the folder containing PDF/DOC files: ")
    output_folder = input("Enter the path to the output folder for TXT files: ")
    process_folder(input_folder, output_folder)
