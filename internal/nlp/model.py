import spacy
from spacy.training.example import Example
from spacy.util import minibatch, compounding
import random
import json

def train():
    nlp = spacy.load("en_core_web_sm")

    # Загрузка тренировочных данных
    with open("train_data.json", "r") as f:
        train_data = json.load(f)

    # Преобразование данных в формат spaCy
    formatted_data = []
    for entry in train_data:
        text = entry["question"]
        entities = []
        for keyword in entry["keywords"]:
            start = text.lower().find(keyword.lower())
            if start != -1:
                end = start + len(keyword)
                entities.append((start, end, keyword))
        formatted_data.append((text, {"entities": entities}))

    # Обучение модели
    optimizer = nlp.begin_training()
    n_iter = 50  # Количество эпох

    print("Starting training...")
    for itn in range(n_iter):
        random.shuffle(formatted_data)
        losses = {}

        batches = minibatch(formatted_data, size=compounding(4.0, 32.0, 1.001))
        for batch in batches:
            examples = []
            for text, annotations in batch:
                doc = nlp.make_doc(text)
                examples.append(Example.from_dict(doc, annotations))
            
            nlp.update(examples, drop=0.5, losses=losses)

        print(f"Iteration {itn + 1} | Losses: {losses}")

    output_dir = "model"
    nlp.to_disk(output_dir)
    print(f"Model saved to {output_dir}")
