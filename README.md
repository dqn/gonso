# go-nso

Array.from(document.querySelectorAll('#authorize-switch-approval-link')).map((it) => it.href.split('&')[1].replace('session_token_code=', ''));
