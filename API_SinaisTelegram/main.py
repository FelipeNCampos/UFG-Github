import time

from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

from featchtelegram import send_message

driver = webdriver.Chrome()


def confereSinal(ultimo, atual, amplitude):
    ultimo = int(ultimo)
    atual = int(atual)
    argumentos = [0, 32, 15, 19, 4, 21, 2, 25, 17, 34, 6, 27, 13, 36, 11, 30, 8, 23, 10, 5, 24, 16, 33, 1, 20, 14, 31,
                  9, 22, 18, 29, 7, 28, 12, 35, 3, 26]
    ianterior = argumentos.index(ultimo)
    grupo = []

    for i in range((amplitude * 2) + 1):
        grupo.append(argumentos[(ianterior - amplitude + i) % 37])

    if (atual in grupo):
        return True
    else:
        return False

def selectXP(path):
    try:
        # Espera explícita até que o elemento seja visível
        element = WebDriverWait(driver, 30).until(
            EC.presence_of_element_located((By.XPATH, path))
        )
        return element
    except Exception as e:
        print(f"Erro ao localizar elemento com XPath '{path}': {e}")
        return None
acaso1 = -1
acaso2 = -1
acaso3 = -1
acaso4 = -1
acaso5 = -1
acaso6 = -1
acaso7 = -1
acaso8 = -1
sequencia = 0
ultimarsp = False
def trackroleta():
    global acaso1,acaso2,acaso3,acaso4,acaso5,acaso6,acaso7,acaso8,sequencia,ultimarsp

    caso1 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[1]").text
    caso2 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[2]").text
    caso3 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[3]").text
    caso4 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[4]").text
    caso5 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[5]").text
    caso6 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[6]").text
    caso7 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[7]").text
    caso8 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[8]").text
    if[caso1,caso2,caso3,caso4,caso5,caso6,caso7,caso8]!=[acaso1,acaso2,acaso3,acaso4,acaso5,acaso6,acaso7,acaso8]:
        acaso1 = caso1
        acaso2 = caso2
        acaso3 = caso3
        acaso4 = caso4
        acaso5 = caso5
        acaso6 = caso6
        acaso7 = caso7
        acaso8 = caso8
        resp = confereSinal(acaso1, acaso2, 7)
        if resp == ultimarsp:
            sequencia += 1
        else:
            sequencia = 0
        if resp == True:
            print("GREEN", end="")
            message: str = "GREEN"
            if sequencia != 0:
                print(" - Sequencia de", sequencia + 1, "GREENs")
                message += f" - Sequencia de {sequencia+1} GREENs"
            else:
                print("")
            send_message('7862862108:AAEARORzD7rLcyl0MflNTC_R8OTnI0g2Oy0','6545737367',message)

        else:
            print("RED", end="")
            message: str = "RED"
            if sequencia != 0:
                print(" - Sequencia de", sequencia + 1, "REDs")
                message += f" - Sequencia de {sequencia+1} REDs"
            else:
                print("")
            send_message('7862862108:AAEARORzD7rLcyl0MflNTC_R8OTnI0g2Oy0', '6545737367', message)
        ultimarsp = resp

    else:
        print('Aguardando dealer')




def main():
    driver.get("https://casino.betfair.com/pt-br/")
    time.sleep(5)
    global acaso1,acaso2,acaso3,acaso4,acaso5,acaso6,acaso7,acaso8
    acaso1 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[1]").text
    acaso2 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[2]").text
    acaso3 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[3]").text
    acaso4 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[4]").text
    acaso5 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[5]").text
    acaso6 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[6]").text
    acaso7 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[7]").text
    acaso8 = selectXP("/html/body/div[2]/div/div[2]/div[2]/div/div[2]/div/div[1]/div/div[1]/div[1]/div/span[8]").text

    while True:
        trackroleta()
        time.sleep(1.5)


if __name__ == '__main__':
    main()