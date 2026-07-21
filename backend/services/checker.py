# services/checker.py

def normalize_answer(answer: str) -> str:
    return answer.strip().lower()


def check_answer(task_type: str, correct: str, user_answer: str) -> bool:
    user = normalize_answer(user_answer)
    correct_norm = normalize_answer(correct)

    if task_type == "choice":
        return user == correct_norm

    elif task_type == "number":
        try:
            user_num = float(user)
            correct_num = float(correct_norm)
            return abs(user_num - correct_num) <= 0.01
        except ValueError:
            return False

    elif task_type == "string":
        return user == correct_norm

    elif task_type == "multi":
        user_set = set(user.split(","))
        correct_set = set(correct_norm.split(","))
        return user_set == correct_set

    return False