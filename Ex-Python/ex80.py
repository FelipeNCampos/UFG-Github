l=[]
for c in range(0,5):
    n=int(input('Digite um numero a ser adicionado : '))
    if c==0 or n>l[-1]:
        l.append(n)
        print(f'O valor {n} foi adicionado no final da lista')
    else :
        p = 0
        while p<len(l):
            if n <=l[p]:
                l.insert(p,n)
                break
            p+=1
        print('Valor adicionado na posicao {}'.format(p))
print(l)