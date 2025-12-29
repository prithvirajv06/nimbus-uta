package internal

func GetWelcomeTemplate(userName string, unsubscribeURL string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Nimbus UTA</title>
    <style>
        body { font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; margin: 0; padding: 0; background-color: #f4f7f9; color: #333; }
        .container { max-width: 600px; margin: 20px auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; border: 1px solid #e1e8ed; }
        .header { background-color: #0f172a; padding: 40px 20px; text-align: center; }
        .header img { max-width: 180px; height: auto; }
        .content { padding: 40px 30px; line-height: 1.6; }
        .content h1 { color: #1e293b; font-size: 24px; margin-top: 0; }
        .cta-button { display: inline-block; padding: 14px 30px; background-color: #3b82f6; color: #ffffff !important; text-decoration: none; border-radius: 5px; font-weight: bold; margin-top: 20px; }
        .footer { background-color: #f8fafc; padding: 20px; text-align: center; font-size: 12px; color: #64748b; }
        .social-links { margin-bottom: 10px; }
        .social-links a { margin: 0 10px; text-decoration: none; color: #3b82f6; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
			<h1>Nimbus UTA</h1>
        </div>

        <div class="content">
            <h1>Welcome to the future, ` + userName + `!</h1>
            <p>We're thrilled to have you join our community. Nimbus UTA is built to provide you with seamless, high-performance solutions tailored to your needs.</p>
            
            <p><strong>To get started, we recommend these three steps:</strong></p>
            <ul>
                <li>Complete your profile setup.</li>
                <li>Explore your new personalized dashboard.</li>
                <li>Join our community discord or forum.</li>
            </ul>

            <a href="` + unsubscribeURL + `" class="cta-button">Launch Your Dashboard</a>
            
            <p style="margin-top: 30px;">If you have any questions, simply reply to this email. Our team is always here to help!</p>
            <p>Best regards,<br>The Nimbus UTA Team</p>
        </div>

        <div class="footer">
            <div class="social-links">
                <a href="#">LinkedIn</a> | <a href="#">Twitter</a> | <a href="#">Instagram</a>
            </div>
            <p>&copy; 2025 Nimbus UTA. All rights reserved.<br>
            123 Tech Plaza, Innovation Way, CA 94043</p>
            <p><a href="{{.UnsubscribeURL}}" style="color: #64748b;">Unsubscribe</a></p>
        </div>
    </div>
</body>
</html>`
}
