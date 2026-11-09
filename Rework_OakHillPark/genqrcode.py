import qrcode

def gerar_qrcode(link, nome_arquivo="Readings.png"):
    """
    Gera um QR Code a partir de um link e o salva como uma imagem.

    :param link: URL ou texto que será convertido em QR Code
    :param nome_arquivo: Nome do arquivo de saída (padrão: qrcode.png)
    """
    try:
        # Gera o QR Code
        qr = qrcode.QRCode(
            version=1,  # Tamanho do QR Code (1 = menor, até 40 = maior)
            error_correction=qrcode.constants.ERROR_CORRECT_L,  # Nível de correção de erros
            box_size=10,  # Tamanho de cada "caixa" no QR Code
            border=4,  # Tamanho da borda (mínimo é 4)
        )
        qr.add_data(link)
        qr.make(fit=True)

        # Cria a imagem do QR Code
        img = qr.make_image(fill_color="black", back_color="white")
        
        # Salva a imagem em um arquivo
        img.save(nome_arquivo)
        print(f"QR Code salvo como {nome_arquivo}")
    except Exception as e:
        print(f"Erro ao gerar QR Code: {e}")

# Exemplo de uso
if __name__ == "__main__":
    link = input("Digite o link para gerar o QR Code: ").strip()
    if link:
        gerar_qrcode(link)
    else:
        print("Nenhum link foi fornecido. Tente novamente.")
