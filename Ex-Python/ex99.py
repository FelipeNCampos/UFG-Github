def maior(x):
    print(f'Dentre os {len(x)} numeros gravados o maior foi {max(x)}')


n=[]
while True:
    n.append(input('Digite um numero a ser gravado :'))
    r=str(input('Deseja continuar (S/N) ? ')).strip().upper()
    if r not in 'SN':
        r=str(input('Resposta invalida, deseja continuar (S/N?'))
    if r=='N':
        break
maior(n)