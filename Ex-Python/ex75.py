t = (int(input('Digite um numero : ')),int(input('Digite outro numero : ')),int(input('Digite outro numero : ')),int(input('Digite outro numero : ')))
qnt9 = t.count(9)
print('O numero nove apareceu {} vezes '.format(qnt9))
id3 = t.index(3)
print('A primeira recorrencia do numero 3 foi no indice {}'.format(int (id3)+1))      
cntp = 0
for c in range(0,4):
    if t[c]%2==0 :
        cntp += 1
print('Foram digitados {} numeros pares '.format(cntp))