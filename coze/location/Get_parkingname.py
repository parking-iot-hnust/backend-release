from runtime import Args
from typings.Get_parkingname.Get_parkingname import Input, Output

def handler(args: Args[Input]) -> Output:
    input = args.input
    logger = args.logger

    # 假设 input.parking_lots 是一个包含 CustomNamespace 对象的列表
    parking_lots = input.parking_lots

    # 检查停车场数量是否超过5个
    if len(parking_lots) > 5:
        raise ValueError("最多只能处理5个停车场名字")

    # 初始化输出字典，所有名称都设为None
    output_names = {
        "name1": None,
        "name2": None,
        "name3": None,
        "name4": None,
        "name5": None
    }

    # 将停车场名称赋值到输出字典中
    for i, parking_lot in enumerate(parking_lots[:5]):  # 只处理前5个停车场
        # 使用点符号访问 CustomNamespace 对象的 name 属性
        if hasattr(parking_lot, 'name'):
            output_names[f"name{i+1}"] = getattr(parking_lot, 'name')
        else:
            logger.warning(f"停车场信息缺少 'name' 属性: {parking_lot}")

    return output_names