import cv2
import numpy as np

def remover_fundo_cor(imagem_entrada, imagem_saida):
    # Carregar a imagem
    img = cv2.imread(imagem_entrada)

    # Converter para espaço de cores HSV
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)

    # Definir limites para o fundo que deseja remover (ex: fundo verde)
    lower_bound = np.array([35, 55, 55])  # Ajuste os valores de acordo com o fundo
    upper_bound = np.array([85, 255, 255])

    # Criar máscara
    mask = cv2.inRange(hsv, lower_bound, upper_bound)
    mask_inv = cv2.bitwise_not(mask)

    # Remover fundo
    img_fg = cv2.bitwise_and(img, img, mask=mask_inv)

    # Salvar imagem com fundo transparente
    alpha = np.zeros_like(mask_inv, dtype=np.uint8)
    result = cv2.merge((img_fg, alpha))
    cv2.imwrite(imagem_saida, result)

# Exemplo de uso
remover_fundo_cor('image.png', 'saida.png')
