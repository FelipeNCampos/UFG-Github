aln = {}
aln['Nome'] = input('Nome : ').title()
aln['Media'] = int(input('Media : '))
if aln['Media']<7:
    aln['Situação'] = 'Reprovado'
if aln['Media'] >=7:
    aln['Situação'] = 'Aprovado'
for c in aln:
    print(f'{c} é igual à {aln.get(c)}')