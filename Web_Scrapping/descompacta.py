import os
import zipfile

# Defina o caminho da pasta
caminho = 'C:\\Users\\felip\\Desktop\\docs'

# Liste todos os arquivos e diretórios na pasta
nomes = os.listdir(caminho)

# Imprima os nomes
for nome in nomes:
    with zipfile.ZipFile(f'C:\\Users\\felip\\Desktop\\docs\\{nome}', 'r') as zip_ref:
        zip_ref.extractall('C:\\Users\\felip\\Desktop\\Unzip')
