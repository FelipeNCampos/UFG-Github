#include<bits/stdc++.h>
using namespace std;

int main(){
    int a;
    int count = 0;
    string linha, temp;
    cin >> a;
    
    cin.ignore();
    for(int i = 0 ;i<=a;i++){
        getline(cin, temp);
        linha+=temp;
    }
    for(int c = 0 ;c<=linha.size()-1;c++){
        if(linha.at(c)=='{'){
            count++;
        }
        else if(linha.at(c)=='}'){
            count--;
        }
        if (count<0){
            break;
        }

    }
    if (count==0){
        cout<<'S'<<'\n';
    }else{
        cout<<'N'<<'\n';
    }
    return 0;
}
