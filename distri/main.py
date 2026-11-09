from selenium.webdriver.common.by import By
from selenium.webdriver import ActionChains
from selenium import webdriver
import time

driver = webdriver.Chrome()
valor_aposta = 0.1
times = 3

def seletorX(link):
    return driver.find_element(By.XPATH, link)

def marcar(op):
    # Marca a aposta múltipla (1x0, 2x0 ou 3x0, por exemplo) com base no op (option) fornecido

    if op==1:
        #marcar 1/0
        local = seletorX('//*[@id="main-view"]/ng-component/div/ms-option-group-list/div[2]/ms-option-panel[1]/div/ms-regular-group/ms-regular-option-group/div/ms-option[1]/ms-event-pick/div/div[1]')
    
    elif op==2:
        #marcar 0/1
        local = seletorX('//*[@id="main-view"]/ng-component/div/ms-option-group-list/div[2]/ms-option-panel[1]/div/ms-regular-group/ms-regular-option-group/div/ms-option[2]/ms-event-pick/div/div[1]')
        
    elif op==3:
        #marcar 2/1
        local = seletorX('//*[@id="main-view"]/ng-component/div/ms-option-group-list/div[2]/ms-option-panel[1]/div/ms-regular-group/ms-regular-option-group/div/ms-option[5]/ms-event-pick/div/div[1]')

    elif op==4:
        #marcar 1/2
        local = seletorX('//*[@id="main-view"]/ng-component/div/ms-option-group-list/div[2]/ms-option-panel[1]/div/ms-regular-group/ms-regular-option-group/div/ms-option[6]/ms-event-pick/div/div[1]')
        
    elif op==5:
        #marcar 3/2
        #mostrar mais
        local = driver.find_element(By.XPATH, '//*[@id="main-view"]/ng-component/div/ms-option-group-list/div[2]/ms-option-panel[1]/div/ms-regular-group/ms-regular-option-group/ms-option-panel-bottom-action/div/span[1]')
        local.click()
        
        local = seletorX('//*[@id="main-view"]/ng-component/div/ms-option-group-list/div[2]/ms-option-panel[1]/div/ms-regular-group/ms-regular-option-group/div/ms-option[7]/ms-event-pick/div/div[1]')
        actions = ActionChains(driver)
        actions.move_to_element(local).perform()
        time.sleep(1)
        
    elif op==6:
        #marcar 2/3
        #mostrar mais
        local = driver.find_element(By.XPATH, '//*[@id="main-view"]/ng-component/div/ms-option-group-list/div[2]/ms-option-panel[1]/div/ms-regular-group/ms-regular-option-group/ms-option-panel-bottom-action/div/span[1]')
        local.click()
        
        local = seletorX('//*[@id="main-view"]/ng-component/div/ms-option-group-list/div[2]/ms-option-panel[1]/div/ms-regular-group/ms-regular-option-group/div/ms-option[8]/ms-event-pick/div/div[1]')
        actions = ActionChains(driver)
        actions.move_to_element(local).perform()
        time.sleep(1)

    local.click()
    time.sleep(2)
    driver.get("https://sports.sportingbet.com/pt-br/sports/favoritos/pr%C3%B3ximos")
    time.sleep(5)
        

# Encontrar o botão que redireciona aos placares de aposta
def irPlacar():
    time.sleep(2)
    for c in range(1,11):  
        campo = seletorX('//*[@id="main-view"]/ng-component/ms-header/ms-dropdown-tab-bar/ul/li['+str(c)+']')
        if 'Placar' in campo.text or 'Mais' in campo.text:
            campo.click()
            time.sleep(0.5)
            if campo.text=='Mais':
                for j in range(1,11):
                    campo2 = seletorX('/html/body/popper-content/div/div[1]/div/div['+str(j)+']')
                    if 'Placar' in campo2.text:
                        campo2.click()
                        return
            return        

def login():#login
    login = "txungopvh@gmail.com"
    senha = "@15202230Br"
    driver.get('https://sports.sportingbet.com/pt-br/sports/favoritos/pr%C3%B3ximos')
    
    time.sleep(3)
    seletorX('//*[@id="userId"]').send_keys(login)
    time.sleep(0.1)
    seletorX('//*[@id="password"]/input').send_keys(senha)
    
    seletorX('//*[@id="login"]/form/fieldset/section/div/button').click()
    
    time.sleep(10)
    
    seletorX('//*[@id="onetrust-accept-btn-handler"]').click()


login()

def isInverted():
    # Retorna True se o time visitante for favoritado (apostaremos nele), se False apostaremos no time da casa
    return (seletorX('//*[@id="main-view"]/ng-component/ms-header/ms-header-content/ms-scoreboard/ms-prematch-scoreboard/div/div[3]/ms-scoreboard-participant/ms-favourite-participant-toggle/ms-favourite-toggle/i').get_attribute('class'))[-8:]=='selected'

for tabela in range(1, 4): # Pra cada tabela de N jogos (3 <= N <= 8)
    for aposta in range(1, 10): # Pra cada jogo da Tabela
        time.sleep(1)
        
        #coluna 1
        local = seletorX('//*[@id="main-view"]/ms-favourites-dashboard/div/ms-grid/div/ms-event-group/ms-event[1]/div/a/ms-event-detail/ms-event-name/ms-inline-tooltip/div/div[2]/div')
        actions = ActionChains(driver)
        actions.move_to_element(local).perform()
        time.sleep(1)
        local.click()
        irPlacar()
        
        if aposta<=3:
            if isInverted():
                marcar(2)
            else:
                marcar(1)
                
        elif aposta>=4 and aposta<=6:
            if isInverted():
                marcar(4)
            else:
                marcar(3)

        elif aposta>=7 and aposta<=9:
            if isInverted():
                marcar(6)
            else:
                marcar(5)

        #coluna 2
        
        seletorX('//*[@id="main-view"]/ms-favourites-dashboard/div/ms-grid/div/ms-event-group/ms-event[2]/div/a/ms-event-detail/ms-event-name/ms-inline-tooltip/div/div[2]/div').click()
        irPlacar()
        if aposta%3==1:
            
            if isInverted():
                marcar(2)# marcar invertido
            else:
                marcar(1)# normal
                
        elif aposta%3==2:
            
            if isInverted():
                marcar(4)# marcar invertido
            else:
                marcar(3)# normal
                
        elif aposta%3==0:
            if isInverted():
                marcar(6)# marcar invertido
            else:
                marcar(5)# normal
            
        
        #coluna 3 iterar para proximas
        
        for k in range(3,times+1):
            local = seletorX('//*[@id="main-view"]/ms-favourites-dashboard/div/ms-grid/div/ms-event-group/ms-event['+str(k)+']/div/a/ms-event-detail/ms-event-name/ms-inline-tooltip/div/div[2]/div')
            actions = ActionChains(driver)
            actions.move_to_element(local).perform()
            time.sleep(1)
            local.click()
            irPlacar()
            if tabela==1:
                if isInverted():
                    marcar(2)
                else:
                    marcar(1)
                    
            elif tabela==2:
                if isInverted():
                    marcar(4)
                else:
                    marcar(3)
                    
            elif tabela==3:
                if isInverted():
                    marcar(6)
                else:
                    marcar(5)

        # Compra do ticket

        time.sleep(4)
        preco = seletorX('//*[@id="main-content"]/ms-main/div[1]/ng-scrollbar[2]/div/div/div/div/ms-widget-column/ms-widget-slot/ms-bet-column/bs-betslip/div/bs-betslip-edit-state/div[2]/bs-digital-summary/div[1]/bs-digital-summary-stake/ms-stake/div/ms-stake-input/div/input')
        actions = ActionChains(driver)
        actions.move_to_element(preco).perform()
        time.sleep(1)
        preco.clear()
        preco.clear()
        preco.send_keys(valor_aposta)
        time.sleep(1)


        botaoComprar = seletorX('//span[@class="ds-btn-text"]') # Seleção do botão de compra
        actions = ActionChains(driver) # Rolar até o elemento
        actions.move_to_element(botaoComprar).perform()
        time.sleep(1)
        botaoComprar.click()
        time.sleep(5)
        seletorX('//span[@class="ds-btn-text"]').click()
        driver.get("https://sports.sportingbet.com/pt-br/sports/favoritos/pr%C3%B3ximos")
        time.sleep(10)
print('Jogos realizados com sucesso!')