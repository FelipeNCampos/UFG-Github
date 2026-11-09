
while True:
    f =str(input('Digite um frase a ter os parenteses analisados :')).strip()
    a=b=r=0
    for x,y in enumerate(f):
        if y in '(':
            a+=1
        if y in ')':
            b+=1
        if b>a:
            break
    if a-b!=0:
        break
    print('Frase\033[032m correta\033[m')
    r = str(input('Deseja continuar (S/N) ? ')).upper().strip()
    if r!='S':
        break
if a-b>0:
    print('Frase\033[031m errada \033[m ,voce nao fechou os parenteses!')
if a-b<0:
    print('Frase\033[031m errada\033[m ,voce nao abriu os parenteses primeiro!')