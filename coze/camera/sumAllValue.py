from runtime import Args
from typings.sumAllValue.sumAllValue import Input, Output

def handler(args: Args[Input])->Output:
    input = args.input
    items = input.Result.items
    sum = 0
    for item in items:
        sum += int(item.value)
    return {"total": str(sum)}