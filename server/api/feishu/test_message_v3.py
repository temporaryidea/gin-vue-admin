import asyncio
import json
import logging
from feishu_message_sender import FeishuMessageSender


logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


async def test_send_message():
    sender = FeishuMessageSender(
        app_id="cli_a29e059c17fa900d",
        app_secret="ugExcDYRRfch5TFsMvaQJfoQfdgrqIgV"
    )
    email = "ruansichun@bytedance.com"
    message = json.dumps({
        "header": {
            "template": "pink",
            "title": {
                "content": "出行预警：您的航班时间已变更",
                "tag": "plain_text"
            }
        },
        "elements": [
                {
                    "tag": "div",
                    "text": {
                        "content": "航班信息：明日航空 MR8073",
                        "tag": "lark_md"
                    }
                },
                {
                    "tag": "div",
                    "text": {
                        "content": "03月01日周五",
                        "tag": "lark_md"
                    }
                },
                {
                    "tag": "hr"
                },
                {
                    "tag": "div",
                    "text": {
                        "content": (
                            "<div style='display: flex; justify-content: space-between;'>"
                            "<div style='flex: 1; text-align: center;'>"
                            "**成都双流 T2**\n"
                            "预计起飞\n"
                            "<font color='red'>21:25</font>\n"
                            "(<font color='grey'>原计划 17:25</font>)\n\n"
                            "<div style='background: #f0f0f0;"
                            " padding: 8px; border-radius: 4px;'>"
                            "值机柜台：J,K</div>"
                            "</div>"
                            "<div style='flex: 1; text-align: center;'>"
                            "✈️\n━ ━ ━ ━ ━\n"
                            "2小时30分\n\n"
                            "<div style='background: #f0f0f0;"
                            " padding: 8px; border-radius: 4px;'>"
                            "登机口：23</div>"
                            "</div>"
                            "<div style='flex: 1; text-align: center;'>"
                            "**杭州萧山 T3**\n"
                            "预计到达\n"
                            "<font color='red'>23:55</font>\n"
                            "(<font color='grey'>原计划 19:55</font>)\n\n"
                            "<div style='background: #f0f0f0;"
                            " padding: 8px; border-radius: 4px;'>"
                            "行李转盘：B5</div>"
                            "</div>"
                            "</div>"
                        ),
                        "tag": "lark_md"
                    }
                },
                {
                    "tag": "action",
                    "actions": [
                        {
                            "tag": "button",
                            "text": {
                                "content": "我要退改签",
                                "tag": "plain_text"
                            },
                            "type": "danger"
                        }
                    ]
                },
                {
                    "tag": "note",
                    "elements": [
                        {
                            "tag": "lark_md",
                            "content": (
                                "ℹ️ <font color='grey'>航班延误也不能晚到机场，"
                                "请参考计划飞行时间值机、安检、登机，以免误机</font>"
                            )
                        }
                    ]
                }
            ]
    })

    try:
        success = await sender.send_message(
            email, message, msg_type="interactive")
        if success:
            logger.info(f"Successfully sent message to {email}")
        else:
            logger.error(f"Failed to send message to {email}")
        return success
    except Exception as e:
        logger.error(f"Error sending message: {str(e)}")
        return False

if __name__ == "__main__":
    asyncio.run(test_send_message())
