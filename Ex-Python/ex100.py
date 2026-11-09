from random import randint
l=[]
from time import sleep

def rnd(l):
    for c in range(0,5):
        l.append(randint(0,10))
        print(l[c], end='')
        if c!=4:
            print(',', end='')
        sleep(0.5)
    print()
    return(l)


def somap(p):
    sum=0
    for x,c in enumerate(p):
        if c%2==0:
            sum+=c
    print(f'A soma dos valores pares sorteados foi {sum}')



rnd(l)
somap(l)