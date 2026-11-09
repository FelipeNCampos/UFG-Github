c = ('Atletico','Flamengo','Palmeiras','Fortaleza','Corinthians','Red Bull Bragantino','Fluminense','America','Atletico','Santos','Ceara','Internacional','Sao Paulo','Athletico Paranaense','Cuiaba','Juentude','Gremio','Bahia','Sport','Chapecoense')
print('A=Os primeiros colocados do brasileirao!')
print('B=Os ultimos 4 colocados do brasileirao!')
print('C=Os times do brasileirao ordenados em ordem alfabetica')
print('D=Em qual posicao esta o time da chapecoense')
r = str(input('Digite a sua escolha : ')).upper().strip()
if r == 'A':
    print('Os primeiros colocados 5 colocados sao ',c[0:5])
elif r== 'B':
    print('Os ultimos 4 colocados do brasileirao sao : {}'.format(c[-4:]))
elif r=='C':
    print('Os times do brasileirao em ordem alfabetica sao : {}'.format(sorted(c)))
elif r=='D':
    print('O time da chapecoense esta na posicao : {}'.format(int(c.index('Chapecoense')+1)))
