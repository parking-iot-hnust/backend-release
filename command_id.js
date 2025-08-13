async function debugGetShortcuts() {
    const token = 'coze-token';
    const botId = "bot-id";
    
    // 1. 测试 Token
    console.log('=== 测试 Token 有效性 ===');
    try {
        const tokenTest = await fetch('https://api.coze.cn/v1/bots', {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const tokenResult = await tokenTest.json();
        console.log('Token 测试:', tokenResult.code === 0 ? '✓ 有效' : '✗ 无效');
    } catch (e) {
        console.log('Token 测试失败:', e.message);
    }
    
    // 2. 尝试获取智能体基本信息
    console.log('\n=== 获取智能体基本信息 ===');
    try {
        const botInfo = await fetch(`https://api.coze.cn/v1/bot/get?bot_id=${botId}`, {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const botResult = await botInfo.json();
        console.log('智能体基本信息:', JSON.stringify(botResult, null, 2));
    } catch (e) {
        console.log('获取基本信息失败:', e.message);
    }
    
    // 3. 尝试获取在线配置
    console.log('\n=== 获取智能体在线配置 ===');
    try {
        const onlineInfo = await fetch(`https://api.coze.cn/v1/bot/get_online_info?bot_id=${botId}`, {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const onlineResult = await onlineInfo.json();
        console.log('在线配置结果:', JSON.stringify(onlineResult, null, 2));
        
        if (onlineResult.data?.shortcuts) {
            onlineResult.data.shortcuts.forEach(shortcut => {
                console.log(`快捷指令ID: ${shortcut.id}, 名称: ${shortcut.name}`);
            });
        }
    } catch (e) {
        console.log('获取在线配置失败:', e.message);
    }
}

debugGetShortcuts();