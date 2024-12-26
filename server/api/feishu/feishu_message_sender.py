import asyncio
import aiohttp
import json
import time
from typing import List, Dict, Union
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class FeishuMessageSender:
    def __init__(self, app_id: str, app_secret: str):
        self.app_id = app_id
        self.app_secret = app_secret
        self.base_url = "https://open.feishu.cn/open-apis"
        self.tenant_access_token = None
        self.rate_limiter = asyncio.Semaphore(45)  # 50 QPS limit, using 45 to be safe
        
    async def get_tenant_access_token(self) -> str:
        """Get tenant access token for authentication"""
        url = f"{self.base_url}/auth/v3/tenant_access_token/internal"
        headers = {"Content-Type": "application/json"}
        data = {
            "app_id": self.app_id,
            "app_secret": self.app_secret
        }
        
        async with aiohttp.ClientSession() as session:
            async with session.post(url, json=data, headers=headers) as response:
                result = await response.json()
                if result.get("code") == 0:
                    return result.get("tenant_access_token")
                raise Exception(f"Failed to get tenant access token: {result}")

    async def get_user_open_id(self, email: str) -> str:
        """Convert email to open_id"""
        if not self.tenant_access_token:
            self.tenant_access_token = await self.get_tenant_access_token()
            
        url = f"{self.base_url}/contact/v3/users/batch_get_id"
        headers = {
            "Authorization": f"Bearer {self.tenant_access_token}",
            "Content-Type": "application/json; charset=utf-8"
        }
        data = {
            "emails": [email],
            "include_resigned": False
        }
        
        async with aiohttp.ClientSession() as session:
            async with session.post(url, json=data, headers=headers) as response:
                result = await response.json()
                if result.get("code") == 0:
                    user_list = result.get("data", {}).get("user_list", [])
                    if user_list and len(user_list) > 0:
                        user = user_list[0]
                        if user.get("status", {}).get("is_activated", False):
                            return user.get("user_id")  # Will be converted to open_id based on user_id_type
                raise Exception(f"Failed to get open_id for email {email}: {result}")

    async def send_message(self, email: str, message: Union[str, dict], msg_type: str = "text") -> bool:
        """Send message to user by email"""
        if not self.tenant_access_token:
            self.tenant_access_token = await self.get_tenant_access_token()
            
        url = f"{self.base_url}/im/v1/messages?receive_id_type=email&uuid={int(time.time())}"
        headers = {
            "Authorization": f"Bearer {self.tenant_access_token}",
            "Content-Type": "application/json; charset=utf-8"
        }
        # Prepare content as JSON string as specified by user
        # Handle different message types
        if msg_type == "interactive":
            # For interactive messages, ensure proper card structure
            if isinstance(message, str):
                try:
                    # Try to parse if it's a JSON string
                    content = json.loads(message)
                except json.JSONDecodeError:
                    # If not JSON, treat as raw content
                    content = message
            else:
                # If dict, ensure it's properly structured
                content = message

            # For interactive messages, ensure content is a valid JSON string
            if isinstance(content, str):
                try:
                    # Validate JSON string
                    json.loads(content)
                except json.JSONDecodeError:
                    raise ValueError("Interactive message content must be a valid JSON string")
            else:
                # Convert dict to JSON string
                content = json.dumps(content, ensure_ascii=False)
        else:
            # For text messages, wrap in text object
            content = json.dumps({"text": message}, ensure_ascii=False)

        # Log the actual content for debugging
        logger.info(f"Message content: {content}")

        # Prepare the request data
        data = {
            "receive_id": email,
            "msg_type": msg_type,
            "content": content  # content is already a JSON string
        }
        
        async with self.rate_limiter:  # Implement rate limiting
            async with aiohttp.ClientSession() as session:
                logger.info(f"Request URL: {url}")
                logger.info(f"Request Headers: {headers}")
                logger.info(f"Request Data: {json.dumps(data, ensure_ascii=False, indent=2)}")
                
                async with session.post(url, json=data, headers=headers) as response:
                    result = await response.json()
                    if result.get("code") == 0:
                        return True
                    if result.get("code") == 99991663:  # Token expired
                        self.tenant_access_token = await self.get_tenant_access_token()
                        return await self.send_message(email, message)
                    raise Exception(f"Failed to send message to {email}: {result}")

    async def process_email(self, email: str, message: str, retries: int = 3) -> bool:
        """Process single email with retries"""
        for attempt in range(retries):
            try:
                await self.send_message(email, message)
                logger.info(f"Successfully sent message to {email}")
                return True
            except Exception as e:
                if attempt == retries - 1:
                    logger.error(f"Failed to process {email} after {retries} attempts: {str(e)}")
                    return False
                logger.warning(f"Attempt {attempt + 1} failed for {email}: {str(e)}")
                await asyncio.sleep(1)  # Wait before retry

    async def send_batch_messages(self, emails: List[str], message: str) -> Dict[str, bool]:
        """Send messages to multiple users"""
        results = {}
        total = len(emails)
        
        for i, email in enumerate(emails, 1):
            success = await self.process_email(email, message)
            results[email] = success
            logger.info(f"Progress: {i}/{total} ({(i/total)*100:.1f}%)")
            
        return results

async def main(emails: List[str]):
    # App credentials
    app_id = "cli_a29e059c17fa900d"
    app_secret = "ugExcDYRRfch5TFsMvaQJfoQfdgrqIgV"
    
    sender = FeishuMessageSender(app_id, app_secret)
    message = "devin 测试消息"
    
    results = await sender.send_batch_messages(emails, message)
    
    # Report results
    success_count = sum(1 for v in results.values() if v)
    logger.info(f"\nFinal Results:")
    logger.info(f"Total processed: {len(results)}")
    logger.info(f"Successful: {success_count}")
    logger.info(f"Failed: {len(results) - success_count}")
    
    # Log failed emails
    failed_emails = [email for email, success in results.items() if not success]
    if failed_emails:
        logger.error("Failed emails:")
        for email in failed_emails:
            logger.error(f"- {email}")

if __name__ == "__main__":
    # Example usage:
    test_emails = []  # Add email list here
    asyncio.run(main(test_emails))
