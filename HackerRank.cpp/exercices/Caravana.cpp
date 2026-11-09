#include<bits/stdc++.h>
using namespace std;
const int n = 1000;
int registro[n];
int entrada, add, temp,i ,meio;
int main(){
    ios::sync_with_stdio(false);
    cin.tie(NULL);
    cin>>entrada;
    for(i =0 ; i<entrada;i++){
        cin>>temp;
        add += temp;
        registro[i] = temp;
    }
    meio = add/(i);
    for (i =0;i<entrada;i++){
        cout<<meio-registro[i]<<"\n";
    }

    return 0;
}
