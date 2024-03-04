package emails

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func generateWelcomeHtmlString(feURL string) string {
	html := fmt.Sprintf(`
		<div style="background-color: rgb(255, 255, 255)">
		<table
		align="center"
		width="100%%"
		border="0"
		cellpadding="0"
		cellspacing="0"
		role="presentation"
		style="
			max-width: 37.5em;
			margin-left: auto;
			margin-right: auto;
			padding-bottom: 3rem;
			padding-top: 1.25rem;
		"
		>
		<tbody>
			<tr style="width: 100%%">
			<td>
				<table
				align="center"
				width="100%%"
				border="0"
				cellpadding="0"
				cellspacing="0"
				role="presentation"
				style="margin-left: auto; margin-right: auto"
				>
				<tbody>
					<tr>
					<td>
						<img
						alt="Snip URL Logo"
						height="100"
						src="https://snip-url.dewodt.com/logo-full.png"
						style="
							display: block;
							outline: none;
							border: none;
							text-decoration: none;
							margin-left: auto;
							margin-right: auto;
							border-radius: 0.375rem;
						"
						width="222"
						class="CToWUd a6T"
						data-bit="iit"
						tabindex="0"
						/>
						<div
						class="a6S"
						dir="ltr"
						style="opacity: 0.01; left: 545px; top: 83px"
						>
						<span
							data-is-tooltip-wrapper="true"
							class="a5q"
							jsaction="JIbuQc:.CLIENT"
							><button
							class="VYBDae-JX-I VYBDae-JX-I-ql-ay5-ays CgzRE"
							jscontroller="PIVayb"
							jsaction="click:h5M12e;clickmod:h5M12e; pointerdown:FEiYhc; pointerup:mF5Elf; pointerenter:EX0mI; pointerleave:vpvbp; pointercancel:xyn4sd; contextmenu:xexox;focus:h06R8; blur:zjh6rb;mlnRJb:fLiPzd;"
							data-idom-class="CgzRE"
							jsname="hRZeKc"
							aria-label="Download attachment "
							data-tooltip-enabled="true"
							data-tooltip-id="tt-c5"
							data-tooltip-classes="AZPksf"
							id=""
							jslog="91252; u014N:cOuCgd,Kr2w4b,xr6bB; 4:WyIjbXNnLWY6MTc4NzI0MzYyNTM2NzM5MDk3NCJd; 43:WyJpbWFnZS9qcGVnIl0."
							>
							<span class="OiePBf-zPjgPe VYBDae-JX-UHGRz"></span
							><span
								class="bHC-Q"
								data-unbounded="false"
								jscontroller="LBaJxb"
								jsname="m9ZlFb"
								soy-skip=""
								ssk="6:RWVI5c"
							></span
							><span
								class="VYBDae-JX-ank-Rtc0Jf"
								jsname="S5tZuc"
								aria-hidden="true"
								><span class="bzc-ank" aria-hidden="true"
								><svg
									height="20"
									viewBox="0 -960 960 960"
									width="20"
									focusable="false"
									class="aoH"
								>
									<path
									d="M480-336 288-528l51-51 105 105v-342h72v342l105-105 51 51-192 192ZM263.717-192Q234-192 213-213.15T192-264v-72h72v72h432v-72h72v72q0 29.7-21.162 50.85Q725.676-192 695.96-192H263.717Z"
									></path></svg></span
							></span>
							<div class="VYBDae-JX-ano"></div>
							</button>
							<div
							class="ne2Ple-oshW8e-J9"
							id="tt-c5"
							role="tooltip"
							aria-hidden="true"
							>
							Download
							</div></span
						>
						</div>
					</td>
					</tr>
				</tbody>
				</table>
				<table
				align="center"
				width="100%%"
				border="0"
				cellpadding="0"
				cellspacing="0"
				role="presentation"
				style="text-align: center"
				>
				<tbody>
					<tr>
					<td>
						<h1
						style="
							text-align: left;
							font-size: 1.5rem;
							line-height: 2rem;
						"
						>
						Welcome to
						<span class="il">Snip</span>
						<span class="il">URL</span>
						</h1>
						<p
						style="
							font-size: 1rem;
							line-height: 1.5rem;
							margin: 16px 0;
							margin-top: 1.25rem;
							margin-bottom: 1.25rem;
							text-align: left;
						"
						>
						Thank you for signin up for
						<span class="il">Snip</span>
						<span class="il">URL</span>! We are excited to have you on board. You can now start shortening your URLs and make them easier to share.
						</p>
						<a
						href="%s"
						style="
							border-radius: 0.375rem;
							background-color: #facc15;
							padding-left: 1.5rem;
							padding-right: 1.5rem;
							padding-top: 1rem;
							padding-bottom: 1rem;
							font-size: 0.875rem;
							line-height: 100%%;
							font-weight: 600;
							color: #422006;
							text-decoration: none;
							display: inline-block;
							max-width: 100%%;
							padding: 16px 24px 16px 24px;
						"
						target="_blank"
						><span></span
						><span
							style="
							max-width: 100%%;
							display: inline-block;
							line-height: 120%%;
							"
							>Try Now</span
						><span></span
						></a>
						<hr
						style="
							width: 100%%;
							border: none;
							border-top: 1px solid #eaeaea;
							margin-top: 1.25rem;
							margin-bottom: 1.25rem;
						"
						/>
						<p
						style="
							font-size: 0.875rem;
							line-height: 1.25rem;
							margin: 16px 0;
							margin-top: 0px;
							margin-bottom: 0px;
							text-align: left;
						"
						>
						Copyright © 2024
						<span class="il">Snip</span>
						<span class="il">URL</span>
						</p>
					</td>
					</tr>
				</tbody>
				</table>
			</td>
			</tr>
		</tbody>
		</table>
	</div>
	`, feURL)
	return html
}

func SendWelcomeEmail(email string) (*resend.SendEmailResponse, error) {
	// Get resend secret
	apiKey := os.Getenv("RESEND_API_KEY")
	feURL := os.Getenv("FE_URL")

	// Send email
	client := resend.NewClient(apiKey)
	params := &resend.SendEmailRequest{
		From:    "Snip URL <snip-url@dewodt.com>",
		To:      []string{email},
		Html:    generateWelcomeHtmlString(feURL),
		Subject: "Welcome to Snip URL",
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		return nil, err
	}

	return sent, nil
}
