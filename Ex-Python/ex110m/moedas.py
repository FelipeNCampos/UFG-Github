def aumentar(p,v,form=True):
    if form==False:
        return p*((100+v)/100)
    if form==True:
        f = '{:.2f}'.format(p*((100+v)/100))
        return f


def diminuir(p,v,form=True):
    if form==False:
        return p*((100-v)/100)
    if form==True:
        f = '{:.2f}'.format(p*((100-v)/100))
        return f

def dobro (p,form=True):
    if form==False:
        return p*2
    if form==True:
        f = '{:.2f}'.format(p*2)
        return f


def metade(p,form=True):
    if form==False:
        return p/2
    if form==True:
        f = '{:.2f}'.format(p/2)
        return f
def forma(n):
    return '{:.2f}'.format(n)

def resumo(p,x,y):
    print('---'*20)
    print(f'{'RESUMO DE VALOR':^60}')
    print('---'*20)
    print(f'{'Preço analisado : ':<40}', end='')
    print(f'{forma(p):>20}')
    print(f'O dobro do valor {p:.2f} é =', end='')
    print(f'{dobro(p):>32}')
    print(f'A metade do valor {p:.2f} é =', end='')
    print(f'{metade(p):>31}')
    print(f'O valor {p:.2f} com aumento de {x}% é =', end='')
    print(f'{aumentar(p,x):>22}')
    print(f'O valor {p:.2f} reduzido de {y}% é =', end='')
    print(f'{diminuir(p,y):>25}')
    print('---'*20)