#!/usr/bin/env python3
"""Converte JSON de trilha (dayX/steps) para o modelo atual de Content/Topic.

Saida:
- <output_dir>/contents.json
- <output_dir>/topics.json
- <output_dir>/db_seed.json

Uso:
  python scripts/convert_desafio_to_db_model.py --input caminho/origem.json --plan desafio_iniciante_ia --language pt-br --output-dir tmp/desafio
"""

from __future__ import annotations

import argparse
import json
import re
from pathlib import Path
from typing import Any


def day_key_sort(key: str) -> int:
    match = re.fullmatch(r"day(\d+)", key.strip())
    if not match:
        return 10**9
    return int(match.group(1))


def normalize_whitespace(text: str) -> str:
    return text.replace("\r\n", "\n").replace("\r", "\n").strip()


def first_text_step(steps: list[dict[str, Any]]) -> dict[str, Any] | None:
    for step in steps:
        if step.get("type") == "text":
            return step
    return None


def quizz_options_from_any(options: list[dict[str, Any]]) -> list[dict[str, Any]]:
    out: list[dict[str, Any]] = []
    for opt in options:
        out.append(
            {
                "text": str(opt.get("text", "")).strip() or "Opcao",
                "is_correct": bool(opt.get("isCorrect", opt.get("is_correct", False))),
                "explanation": str(opt.get("explanation", "Sem explicacao")).strip() or "Sem explicacao",
            }
        )
    return out


def match_pairs_from_component(pairs: list[dict[str, Any]]) -> list[dict[str, str]]:
    out: list[dict[str, str]] = []
    for pair in pairs:
        left = str(pair.get("left", "")).strip()
        right = str(pair.get("right", "")).strip()
        if left and right:
            out.append({"text": left, "match": right})
    return out


def as_string_list(values: Any) -> list[str]:
    if not isinstance(values, list):
        return []
    return [str(v).strip() for v in values if str(v).strip()]


def topic_from_step(step: dict[str, Any], order: int, day_id: int) -> dict[str, Any] | None:
    step_type = str(step.get("type", "")).strip()

    # 1) Explicacoes
    if step_type == "text":
        title = str(step.get("title", "Topico")).strip() or "Topico"
        content = normalize_whitespace(str(step.get("content", "")).strip() or "Sem conteudo")
        return {
            "day_id": day_id,
            "order": order,
            "behavior": "explicacao",
            "activity_type": None,
            "payload": {"title": title, "content": content},
        }

    # 2) Quiz nativo
    if step_type == "quiz":
        title = str(step.get("title", "Quiz")).strip() or "Quiz"
        question = normalize_whitespace(str(step.get("question", "Pergunta")).strip() or "Pergunta")
        options = quizz_options_from_any(step.get("options", []))
        if not options:
            options = [
                {
                    "text": "Opcao 1",
                    "is_correct": True,
                    "explanation": "Opcao padrao criada automaticamente",
                }
            ]
        return {
            "day_id": day_id,
            "order": order,
            "behavior": "atividade",
            "activity_type": "Quizz",
            "payload": {"title": title, "question": question, "options": options},
        }

    # 3) Componentes
    if step_type == "component":
        component_name = str(step.get("componentName", "")).strip()
        props = step.get("props", {}) if isinstance(step.get("props"), dict) else {}

        if component_name == "MatchWords":
            title = str(props.get("title", "MatchWords")).strip() or "MatchWords"
            description = normalize_whitespace(
                str(props.get("description", "Associe os pares corretos.")).strip() or "Associe os pares corretos."
            )
            pairs = match_pairs_from_component(props.get("pairs", []))
            if not pairs:
                pairs = [{"text": "Exemplo", "match": "Correspondencia"}]
            return {
                "day_id": day_id,
                "order": order,
                "behavior": "atividade",
                "activity_type": "MatchWords",
                "payload": {"title": title, "description": description, "pairs": pairs},
            }

        if component_name == "FillBlanks":
            title = str(props.get("title", "FillBlanks")).strip() or "FillBlanks"
            sentence = normalize_whitespace(str(props.get("sentence", "Complete as lacunas.")).strip() or "Complete as lacunas.")
            answers = as_string_list(props.get("answers", []))
            answer = " | ".join(answers) if answers else "Resposta"
            options = as_string_list(props.get("options", []))
            if not options:
                options = [answer]
            explanation = normalize_whitespace(str(props.get("explanation", "Resolucao guiada.")).strip() or "Resolucao guiada.")
            return {
                "day_id": day_id,
                "order": order,
                "behavior": "atividade",
                "activity_type": "FillBlanks",
                "payload": {
                    "title": title,
                    "sentence": sentence,
                    "answer": answer,
                    "options": options,
                    "explanation": explanation,
                },
            }

        if component_name in {"Quiz", "Quizz"}:
            title = str(props.get("title", "Quiz")).strip() or "Quiz"
            question = normalize_whitespace(str(props.get("question", "Pergunta")).strip() or "Pergunta")
            options = quizz_options_from_any(props.get("options", []))
            if not options:
                options = [
                    {
                        "text": "Opcao 1",
                        "is_correct": True,
                        "explanation": "Opcao padrao criada automaticamente",
                    }
                ]
            return {
                "day_id": day_id,
                "order": order,
                "behavior": "atividade",
                "activity_type": "Quizz",
                "payload": {"title": title, "question": question, "options": options},
            }

        if component_name == "PromptTrainer":
            title = str(props.get("title", "Prompt Trainer")).strip() or "Prompt Trainer"
            description = normalize_whitespace(
                str(props.get("description", "Pratique um prompt com foco no objetivo.")).strip()
                or "Pratique um prompt com foco no objetivo."
            )
            ai_tool_context = normalize_whitespace(
                str(props.get("aiToolContext", "Contexto de pratica")).strip() or "Contexto de pratica"
            )
            success_criteria_list = as_string_list(props.get("successCriteria", []))
            success_criteria = (
                "\n".join(f"- {item}" for item in success_criteria_list)
                if success_criteria_list
                else "- Definir contexto\n- Definir tarefa\n- Definir formato"
            )
            return {
                "day_id": day_id,
                "order": order,
                "behavior": "atividade",
                "activity_type": "PromptTreiner",
                "payload": {
                    "title": title,
                    "description": description,
                    "AIToolContext": ai_tool_context,
                    "successCriteria": success_criteria,
                },
            }

        # Componente desconhecido => explicacao de fallback
        fallback_title = str(props.get("title") or step.get("title") or component_name or "Componente").strip() or "Componente"
        fallback_content = normalize_whitespace(
            str(props.get("description") or f"Componente {component_name} convertido para explicacao.").strip()
        )
        return {
            "day_id": day_id,
            "order": order,
            "behavior": "explicacao",
            "activity_type": None,
            "payload": {"title": fallback_title, "content": fallback_content},
        }

    # 4) Practical (nao existe no schema atual) => FillBlanks
    if step_type == "practical":
        title = str(step.get("title", "Atividade Pratica")).strip() or "Atividade Pratica"
        question = normalize_whitespace(str(step.get("question", "Ordene os itens.")).strip() or "Ordene os itens.")
        options = as_string_list(step.get("initialWords", []))
        correct_order = as_string_list(step.get("correctOrder", []))
        answer = " | ".join(correct_order) if correct_order else "Resposta"
        explanation = normalize_whitespace(str(step.get("explanation", "Atividade de ordenacao.")).strip() or "Atividade de ordenacao.")
        if not options:
            options = correct_order or ["Opcao 1"]
        return {
            "day_id": day_id,
            "order": order,
            "behavior": "atividade",
            "activity_type": "FillBlanks",
            "payload": {
                "title": title,
                "sentence": question,
                "answer": answer,
                "options": options,
                "explanation": explanation,
            },
        }

    return None


def convert(raw: dict[str, Any], plan: str, language: str) -> dict[str, Any]:
    day_items = []
    for key, value in raw.items():
        if not re.fullmatch(r"day\d+", key):
            continue
        if not isinstance(value, dict):
            continue
        steps = value.get("steps")
        if not isinstance(steps, list):
            continue
        day_items.append((day_key_sort(key), key, steps))

    day_items.sort(key=lambda item: item[0])

    contents: list[dict[str, Any]] = []
    topics: list[dict[str, Any]] = []

    for day_id, _, steps in day_items:
        text_intro = first_text_step(steps)
        if text_intro:
            content_title = str(text_intro.get("title", f"Dia {day_id}")).strip() or f"Dia {day_id}"
            content_body = normalize_whitespace(str(text_intro.get("content", "")).strip() or f"Conteudo do dia {day_id}")
        else:
            content_title = f"Dia {day_id}"
            content_body = f"Conteudo estruturado do dia {day_id}."

        contents.append(
            {
                "plan": plan,
                "language": language,
                "day_id": day_id,
                "content_index": 1,
                "title": content_title,
                "content": content_body,
            }
        )

        order = 1
        for step in steps:
            if not isinstance(step, dict):
                continue
            topic = topic_from_step(step, order=order, day_id=day_id)
            if topic is None:
                continue
            topics.append(topic)
            order += 1

    return {"plan": plan, "language": language, "contents": contents, "topics": topics}


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--input", required=True, help="Arquivo JSON com day1..dayN")
    parser.add_argument("--plan", default="desafio_iniciante_ia")
    parser.add_argument("--language", default="pt-br")
    parser.add_argument("--output-dir", default="tmp/desafio_iniciante_ia")
    args = parser.parse_args()

    input_path = Path(args.input)
    output_dir = Path(args.output_dir)
    output_dir.mkdir(parents=True, exist_ok=True)

    raw = json.loads(input_path.read_text(encoding="utf-8"))
    converted = convert(raw, plan=args.plan, language=args.language)

    contents_path = output_dir / "contents.json"
    topics_path = output_dir / "topics.json"
    merged_path = output_dir / "db_seed.json"

    contents_path.write_text(
        json.dumps(converted["contents"], ensure_ascii=False, indent=2),
        encoding="utf-8",
    )
    topics_path.write_text(
        json.dumps(converted["topics"], ensure_ascii=False, indent=2),
        encoding="utf-8",
    )
    merged_path.write_text(
        json.dumps(converted, ensure_ascii=False, indent=2),
        encoding="utf-8",
    )

    print(f"Conteudos: {len(converted['contents'])}")
    print(f"Topicos: {len(converted['topics'])}")
    print(f"Arquivo completo: {merged_path}")


if __name__ == "__main__":
    main()
