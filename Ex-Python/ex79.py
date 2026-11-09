r='S'
l=[]
while r=='S':
    n=input('Digite um numero a ser gravado :')
    if l.count(n)!=0:
        print('\033[31mValor duplicado\033[m, nao vou adicionar...')
    if l.count(n)==0:
        l.append(n)
        print('Valor adicionado com \033[032msucesso\033[m!')
    r=input('Deseja continuar (S/N)?').upper().strip()
print(f'Os valores digitados foram {l}')