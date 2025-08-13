from runtime import Args
from typings.Parking_charging_system.Parking_charging_system import Input, Output

class Input:
    total_seconds : int
    reg1 : float # 起步价
    reg2 : float # 最大费用
    reg3 : float # 单位费用
    unit : int
    free_time : int

def handler(args: Args[Input])->Output:
    input = args.input
    reg1 = input.reg1
    reg2 = input.reg2
    reg3 = input.reg3
    unit = input.unit
    free_time = input.free_time
    logger = args.logger

    if(input.total_seconds == 0) :
        return {"price" : 100000}
    times = (input.total_seconds + 59) // 60
    res = 0
    if(times <= free_time) :
        res = 0
    else :
        res = reg1
        res += min(reg2, ((times + unit - 1) // unit - 1) * reg3)
    return {"price": res}