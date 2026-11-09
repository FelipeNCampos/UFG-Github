from selenium import webdriver
from selenium.webdriver.common.by import By
import time

# Configurações do Chrome
options = webdriver.ChromeOptions()
options.add_experimental_option('prefs', {
    "download.default_directory": "C:\\Users\\felip\\Desktop\\docs",  # Altere para o seu caminho
    "download.prompt_for_download": False,
    "download.directory_upgrade": True,
    "safebrowsing.enabled": True
})

# Inicia o driver
driver = webdriver.Chrome(options=options)

def seletorX(link):
    return driver.find_element(By.XPATH, link)


try:
    # Acesse a página que contém o botão de download
    driver.get('https://portaldatransparencia.gov.br/download-de-dados/servidores')

    time.sleep(3)
    
    # Pegue todos os elementos que serão utilizados da página
    campoAno = seletorX('//*[@id="links-anos"]')
    campoMes = seletorX('//*[@id="links-meses"]')
    campoTipo = seletorX('//*[@id="links-origens-mes"]')
    
    # Para todos os anos, partindo de 2019 até 2023
    for ano in range(6,1,-1):
        campoAno.click() # Selecione o campo de anos
        
        anoAtual = seletorX(f'//*[@id="links-anos"]/option[{ano}]')
        anoAtual.click() # Selecione o ano

        for mes in range(1,13):
            campoMes.click() # Selecione o campo de meses
                
            mesAtual = seletorX(f'//*[@id="links-meses"]/option[{mes}]')
            mesAtual.click() # Selecione o mês

            campoTipo.click() # Selecione o campo de tipos de planilha
            if ano==6:
                tipo = seletorX('//*[@id="links-origens-mes"]/option[2]')
            else:
                tipo = seletorX('//*[@id="links-origens-mes"]/option[4]')
            
            tipo.click() # Selecione os Honorarios_Jetons


            botaoDownload = seletorX('//*[@id="link"]') # Selecione o botão de download
            botaoDownload.click() # Faça o download da planilha
        
    campoAno.click() # Selecione o campo de anos
    
    anoAtual = seletorX('//*[@id="links-anos"]/option[1]')
    anoAtual.click() # Selecione o ano de 2024

    for mes in range(1,8):
        campoMes.click() # Selecione o campo de meses
                
        mesAtual = seletorX(f'//*[@id="links-meses"]/option[{mes}]')
        mesAtual.click() # Selecione o mes

        campoTipo.click() # Selecione o campo de tipos de planilha
            
        tipo = seletorX('//*[@id="links-origens-mes"]/option[4]')
        tipo.click() # Selecione os Honorarios_Jetons
        
        botaoDownload = seletorX('//*[@id="link"]') # Selecione o botão de download
        botaoDownload.click() # Faça o download da planilha
        
    # Aguardar o download ser concluído (pode ser necessário ajustar o tempo)
    time.sleep(10)  # Aumente se o download demorar mais

finally:
    # Feche o driver
    driver.quit()

