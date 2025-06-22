from pydantic import BaseModel

from tasks.mail_enterprise_task import send_enterprise_email_task


class BossMail(BaseModel):
    to: list[str]
    subject: str
    body: str
    substitutions: dict[str, str] = {}


class EnterpriseMailService:
    @classmethod
    def send_mail(cls, mail: BossMail):
        send_enterprise_email_task.delay(
            to=mail.to, subject=mail.subject, body=mail.body, substitutions=mail.substitutions
        )
