ex = 'Zero','Um','Dois','Três','Quatro','Cinco','Seis','Sete','Oito','Nove','Dez','Onze','Doze','Treze','Quatorze','Quinze','Dezesseis','Dezessete','Dezoito','Dezenove','Vinte'
ep = int(input('Digite um numero de 0 a 20 a ser escrito por extenso : '))
if ep<=20 and ep>=0 :
    for c in range(0,20) :
        if ep == c :
            print(ex[c])
else :
    print('Numero fora da escala ')