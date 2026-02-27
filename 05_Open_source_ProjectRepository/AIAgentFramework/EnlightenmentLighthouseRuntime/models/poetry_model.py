# 诗歌创作模型

def generate_poetry(prompt):
    """生成诗歌"""
    return f"基于提示 '{prompt}'，生成诗歌：
星空下的思绪
如繁星点点
在夜的怀抱中
轻轻摇曳

风穿过指尖
带走了所有的烦恼
留下的
是内心的宁静

每一个瞬间
都是生命的礼物
珍惜当下
便是最好的修行"

if __name__ == "__main__":
    import sys
    if len(sys.argv) > 1:
        prompt = sys.argv[1]
        print(generate_poetry(prompt))