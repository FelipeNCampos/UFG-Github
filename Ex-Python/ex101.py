def voto(nc):
    from datetime import datetime
    n=datetime.now().strftime('%Y')
    if int(n)-nc<16:
        return('Negado')
    if int(n)-nc>16 and int(n)-nc<18:
        return('Opcional')
    if int(n)-nc>18:
        return('Obrigatório')


print('---'*30)
ano = int(input('Digite o ano de nascimento :'))
resp = voto(ano)
print(f'Para a pessoa que nasceu no ano de {ano} o voto é {resp}')