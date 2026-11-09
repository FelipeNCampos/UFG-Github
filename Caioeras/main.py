from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options

import time

# Configurações do Chrome
chrome_options = Options()
chrome_options.add_experimental_option("prefs", {
    "download.default_directory": "C:\\Users\\coton\\Desktop\\caio",  # Diretório para downloads
    "download.prompt_for_download": False,  # Não solicitar confirmação
    "download.directory_upgrade": True,
    "safebrowsing.enabled": True  # Evita bloqueio de downloads perigosos
})

# Inicia o driver
driver = webdriver.Chrome(options=chrome_options)

def seletorX(link):
    return driver.find_element(By.XPATH, link)

#credenciais
login = "ceiras@cahenri.com.br"
senha = 60228228
tempo = 4

try:
    # Acesse a página que contém o botão de download
    driver.get('https://app.orcafascio.com/login/new')
    time.sleep(tempo)
    
    # Pegue todos os elementos que serão utilizados no login
    inputEmailLogin = seletorX('/html/body/div[2]/div/div[2]/div/form/div[1]/div/div[2]/input')
    inputSenhaLogin = seletorX('/html/body/div[2]/div/div[2]/div/form/div[1]/div/div[3]/input')
    botaoEntrarLogin = seletorX('/html/body/div[2]/div/div[2]/div/form/div[2]/button')
    
    #login
    inputEmailLogin.send_keys(login)
    inputSenhaLogin.send_keys(senha)
    botaoEntrarLogin.click()
    time.sleep(tempo)
        
    #pesquisa
    botaoComposicaoMenu = seletorX('/html/body/div[7]/ul/li[3]/a')
    botaoComposicaoMenu.click()
    time.sleep(tempo/3)
    
    botaoPesquisarComposicao = seletorX('/html/body/div[7]/ul/li[3]/ul/li[1]/a')
    botaoPesquisarComposicao.click()
    time.sleep(tempo/3)
    
    for c in range (1,259):
        for index in range (2,32):
            botaoPropriaComposicoes = seletorX('/html/body/div[9]/div[3]/div/div/div[1]/form/div[2]/div[2]/div[2]/label[2]/input')
            botaoPropriaComposicoes.click()
            
            botaoBuscarComposicoes = seletorX('/html/body/div[9]/div[3]/div/div/div[1]/form/div[1]/div[1]/button')
            botaoBuscarComposicoes.click()
            time.sleep(tempo/3)
            
            for page in range (1,c):
                
                if page==1:
                    botaoProximo = seletorX('/html/body/div[9]/div[3]/div/div/div[2]/div[1]/div/div[2]/div[5]/a/img')
                elif page==2:    
                    botaoProximo = seletorX('/html/body/div[9]/div[3]/div/div/div[2]/div[1]/div/div[2]/div[6]/a/img')
                else:
                    botaoProximo = seletorX('/html/body/div[9]/div[3]/div/div/div[2]/div[1]/div/div[2]/div[7]/a/img')
                    
                botaoProximo.click()
                time.sleep(tempo/3)
            
            botaoDoc = seletorX('/html/body/div[9]/div[3]/div/div/div[2]/div[2]/div/div['+str(index)+']/div[2]/div/a')
            botaoDoc.click()
            time.sleep(tempo/5)
            
            if index==2:
                time.sleep(tempo)
                
            botaoFerramenta = seletorX('/html/body/div[4]/div[1]/div[2]/ul/li[2]/a')
            botaoFerramenta.click()
            
            
            botaoExportar = seletorX('/html/body/div[4]/div[1]/div[2]/ul/li[2]/ul/li[2]/a')
            botaoExportar.click()
            
            botaoVoltar = seletorX('/html/body/div[4]/div[3]/ul/li[2]/a')
            botaoVoltar.click()
            time.sleep(tempo/5)

except:
    time.sleep(20)
finally:
    # Feche o driver
    driver.quit()
    
