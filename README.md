Analisis Hugot

   Hugot adalah library Go untuk menjalankan Hugging Face transformer pipelines menggunakan ONNX models. Mendukung:

   Pipelines yang tersedia:
   1. FeatureExtraction - untuk embeddings (RAG, semantic search)
   2. TextClassification - sentiment analysis, dll
   3. TokenClassification - NER (Named Entity Recognition)
   4. ZeroShotClassification - klasifikasi tanpa training
   5. TextGeneration - generate text (Phi, Gemma models)
   6. CrossEncoder - reranking untuk search
   7. ImageClassification - klasifikasi gambar

   Potensi utilisasi di Veridium:

   Use Case                    │ Pipeline            │ Benefit                                     
   ----------------------------+---------------------+---------------------------------------------
   **Local Embeddings**        │ FeatureExtraction   │ Ganti API embeddings dengan local ONNX model
   **Document Classification** │ TextClassification  │ Kategorisasi otomatis file
   **Entity Extraction**       │ TokenClassification │ Extract nama, tanggal, dll dari dokumen
   **Semantic Search**         │ CrossEncoder        │ Reranking search results
   **Local LLM**               │ TextGeneration      │ Alternative ke llama.cpp dengan Phi/Gemma

   Keuntungan:
   •  Pure Go backend untuk deployment simple (tanpa cgo)
   •  ONNX Runtime untuk performance (dengan GPU support)
   •  Bisa fine-tune embedding models
   •  Download model dari Hugging Face otomatis

   Mau saya implementasikan integrasi untuk use case spesifik? Misalnya:
   •  A: Local embeddings dengan FeatureExtraction (ganti embedding API)
   •  B: Document classification untuk auto-categorize files
   •  C: Reranking dengan CrossEncoder untuk better search