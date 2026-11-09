t = ('Desprezar','Suar','chao','forca','pensar')
for c in range(0,5):
    print(f'Na palavra {t[c]} temos as vogais : ', end='')
    if t[c].count('a')!=0 :
        print('a', end='')
    if t[c].count ('e')!=0:
        print(',e', end='')
    if t[c].count ('i')!=0:
        print(',i', end='')
    if t[c].count ('o')!=0:
        print(',o', end='')
    if t[c].count ('u')!=0:
        print(',u', end='')
    print('')
