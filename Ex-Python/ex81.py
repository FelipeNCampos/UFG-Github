t = [(input('Digite um numero a ser gravado : '))]
r = input('Deseja continuar (S/N) ?').upper().strip()
while r=='S':
    t.append(input('Digite outro numero a ser gravado : '))
    r = input('Deseja continuar (S/N) ?').upper().strip()
n1 = int(len(t))-1
n3 = int(len(t))
print(n1)
print(n3)
for c in range(0,n1):
    for k in range(c+1,n3):
        if t[c]<t[k]:
            h=t[c]
            t[c]=t[k]
            t[k]=h
print('Foram digitados um total de {} itens, sendo esses listados de forma decrescente como : {} '.format(len(t),t))
if t.count('5')>0:
    print('Contem o numero 5!')
else:
    print('Nenhum dos numeros gravados é o numero 5')