from playwright.sync_api import sync_playwright
from tqdm import tqdm

def fetch_data():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        page.goto('https://certidao-unificada.cjf.jus.br/#/solicitacao-certidao')
        
        # Exemplo: Aguardar por um seletor que indica que a página está carregada
        page.wait_for_selector(r'/html/body/app-root/app-main/div/div[2]/app-solicitacao-certidao/cjf-panel/div/div[2]/cjf-panel-content/div/p-card[2]/div/div/div/form/div[4]/p-inputmask/input')
        
        # Interagir com elementos
        page.fill(r'/html/body/app-root/app-main/div/div[2]/app-solicitacao-certidao/cjf-panel/div/div[2]/cjf-panel-content/div/p-card[2]/div/div/div/form/div[4]/p-inputmask/input', '04213544109')
        page.click('button[type="submit"]')
        
        # Capturar o conteúdo resultante
        content = page.content()
        
        print(content)  # ou fazer parsing com BeautifulSoup
        browser.close()

tqdm(fetch_data())
